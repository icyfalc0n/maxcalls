package ice

import "github.com/pion/stun/v3"

func parseStunServers(stunServers []string) ([]*stun.URI, error) {
	var result []*stun.URI

	for _, server := range stunServers {
		uri, err := stun.ParseURI(server)
		if err != nil {
			return []*stun.URI{}, err
		}
		result = append(result, uri)
	}

	return result, nil
}

func parseTurnServers(turnServers []string, username string, password string) ([]*stun.URI, error) {
	var result []*stun.URI

	for _, server := range turnServers {
		uri, err := stun.ParseURI(server)
		if err != nil {
			return []*stun.URI{}, err

		}
		uri.Username = username
		uri.Password = password
		result = append(result, uri)
	}

	return result, nil
}
