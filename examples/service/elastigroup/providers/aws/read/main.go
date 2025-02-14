package main

import (
	"context"
	"log"

	"github.com/povils/spotinst-sdk-go/service/elastigroup"
	"github.com/povils/spotinst-sdk-go/service/elastigroup/providers/aws"
	"github.com/povils/spotinst-sdk-go/spotinst"
	"github.com/povils/spotinst-sdk-go/spotinst/session"
	"github.com/povils/spotinst-sdk-go/spotinst/util/stringutil"
)

func main() {
	// All clients require a Session. The Session provides the client with
	// shared configuration such as account and credentials.
	// A Session should be shared where possible to take advantage of
	// configuration and credential caching. See the session package for
	// more information.
	sess := session.New()

	// Create a new instance of the service's client with a Session.
	// Optional spotinst.Config values can also be provided as variadic
	// arguments to the New function. This option allows you to provide
	// service specific configuration.
	svc := elastigroup.New(sess)

	// Create a new context.
	ctx := context.Background()

	// Read group configuration.
	out, err := svc.CloudProviderAWS().Read(ctx, &aws.ReadGroupInput{
		GroupID: spotinst.String("sig-12345"),
	})
	if err != nil {
		log.Fatalf("spotinst: failed to read group: %v", err)
	}

	// Output.
	if out.Group != nil {
		log.Printf("Group %q: %s",
			spotinst.StringValue(out.Group.ID),
			stringutil.Stringify(out.Group))
	}
}
