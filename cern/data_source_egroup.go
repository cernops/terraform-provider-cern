package cern

import (
	"fmt"
	"log"
	"regexp"

	"github.com/go-ldap/ldap"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCernEgroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCernEgroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"members": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func egroup2users(egroups []string, conn *ldap.Conn, recursion bool, processed []string) ([]string, []string) {
	base := "OU=e-groups,OU=Workgroups,DC=cern,DC=ch"
	scope := ldap.ScopeWholeSubtree
	attr := "member"

	var users []string

	for _, egroup := range egroups {
		_, found := Find(processed, egroup)
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
	config := meta.(*Config)

	l, err := ldap.DialURL(config.LdapServer)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer l.Close()

	d.SetId(d.Get("name").(string))

	users, _ := egroup2users([]string{d.Get("name").(string)}, l, true, []string{})
	d.Set("members", users)

	return nil
}
