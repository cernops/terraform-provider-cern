package cern

import (
	"fmt"
	"log"
	"regexp"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCernEgroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCernEgroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name of the e-group to query",
				Required:    true,
			},
			"query_mails": {
				Type:        schema.TypeBool,
				Description: "Flag to specify whether 'mails' should be populated or not",
				Optional:    true,
				Default:     false,
			},
			"members": {
				Type:        schema.TypeList,
				Description: "List of usernames in the e-group",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"mails": {
				Type:        schema.TypeList,
				Description: "List of e-mail addresses of the members of the group",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func getMail(usernames []string, conn *ldap.Conn) []string {
	base := "OU=Users,OU=Organic Units,DC=cern,DC=ch"
	scope := ldap.ScopeWholeSubtree
	attr := "mail"

	var mails []string
	for _, user := range usernames {
		searchRequest := ldap.NewSearchRequest(
			base, scope, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(CN=%s)", user),
			[]string{attr},
			nil,
		)
		sr, err := conn.Search(searchRequest)
		if err != nil {
			log.Fatal(err)
		}

		mails = append(mails, sr.Entries[0].GetAttributeValues(attr)[0])
	}
	return mails
}

func egroup2users(egroups []string, conn *ldap.Conn, recursion bool, processed []string) ([]string, []string) {
	base := "OU=e-groups,OU=Workgroups,DC=cern,DC=ch"
	scope := ldap.ScopeWholeSubtree
	attr := "member"

	var users []string

	for _, egroup := range egroups {
		_, found := find(processed, egroup)
		if found {
			continue
		}
		processed = append(processed, egroup)

		searchRequest := ldap.NewSearchRequest(
			base, scope, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(&(objectClass=group)(CN=%s))", egroup),
			[]string{"member"},
			nil,
		)

		sr, err := conn.Search(searchRequest)
		if err != nil {
			log.Fatal(err)
		}

		userRe := regexp.MustCompile(`CN=(\S+),OU=Users,OU=Organic Units,DC=cern,DC=ch`)
		egroupRe := regexp.MustCompile(`CN=(\S+),OU=e-groups,OU=Workgroups,DC=cern,DC=ch`)

		for _, member := range sr.Entries[0].GetAttributeValues(attr) {
			if userRe.MatchString(member) {
				user := userRe.FindStringSubmatch(member)[1]
				users = append(users, user)
			}
			if egroupRe.MatchString(member) {
				newEgroup := egroupRe.FindStringSubmatch(member)[1]
				expandedUsers, expandedProcessed := egroup2users([]string{newEgroup}, conn, recursion, processed)
				users = append(users, expandedUsers...)
				processed = append(processed, expandedProcessed...)
			}
		}
	}
	return users, processed
}

func dataSourceCernEgroupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config)

	l, err := ldap.DialURL(config.LdapServer)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer l.Close()

	d.SetId(d.Get("name").(string))

	queryMails := d.Get("query_mails").(bool)

	users, _ := egroup2users([]string{d.Get("name").(string)}, l, true, []string{})
	d.Set("members", users)

	if queryMails {
		d.Set("mails", getMail(users, l))
	} else {
		d.Set("mails", []string{})
	}

	return nil
}
