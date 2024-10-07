// credit prometheus alertmanager default template
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Pair is a key/value string pair.
type Pair struct {
	Name, Value string
}

// Pairs is a list of key/value string pairs.
type Pairs []Pair

// Names returns a list of names of the pairs.
func (ps Pairs) Names() []string {
	ns := make([]string, 0, len(ps))
	for _, p := range ps {
		ns = append(ns, p.Name)
	}
	return ns
}

// Values returns a list of values of the pairs.
func (ps Pairs) Values() []string {
	vs := make([]string, 0, len(ps))
	for _, p := range ps {
		vs = append(vs, p.Value)
	}
	return vs
}

func (ps Pairs) String() string {
	b := strings.Builder{}
	for i, p := range ps {
		b.WriteString(p.Name)
		b.WriteRune('=')
		b.WriteString(p.Value)
		if i < len(ps)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

// KV is a set of key/value string pairs.
type KV map[string]string

// SortedPairs returns a sorted list of key/value pairs.
func (kv KV) SortedPairs() Pairs {
	var (
		pairs     = make([]Pair, 0, len(kv))
		keys      = make([]string, 0, len(kv))
		sortStart = 0
	)
	for k := range kv {
		if k == "alertname" {
			keys = append([]string{k}, keys...)
			sortStart = 1
		} else {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys[sortStart:])

	for _, k := range keys {
		pairs = append(pairs, Pair{k, kv[k]})
	}
	return pairs
}

// Remove returns a copy of the key/value set without the given keys.
func (kv KV) Remove(keys []string) KV {
	keySet := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		keySet[k] = struct{}{}
	}

	res := KV{}
	for k, v := range kv {
		if _, ok := keySet[k]; !ok {
			res[k] = v
		}
	}
	return res
}

// Names returns the names of the label names in the LabelSet.
func (kv KV) Names() []string {
	return kv.SortedPairs().Names()
}

// Values returns a list of the values in the LabelSet.
func (kv KV) Values() []string {
	return kv.SortedPairs().Values()
}

func (kv KV) String() string {
	return kv.SortedPairs().String()
}

type Alert struct {
	Status       string `json:"status"`
	Labels       KV     `json:"labels"`
	Annotations  KV     `json:"annotations"`
	StartsAt     string `json:"startsAt"`
	EndsAt       string `json:"endsAt"`
	GeneratorURL string `json:"generatorURL"`
	Fingerprint  string `json:"fingerprint"`
}

type Alerts []Alert

func (as Alerts) Firing() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == "firing" {
			res = append(res, a)
		}
	}
	return res
}

func (as Alerts) Resolved() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == "resolved" {
			res = append(res, a)
		}
	}
	return res
}

type DataPayload struct {
	CustomerName string `json:"-"`
	Receiver     string `json:"receiver"`
	Status       string `json:"status"`
	Alerts       Alerts `json:"alerts"`

	GroupLabels       KV `json:"groupLabels"`
	CommonLabels      KV `json:"commonLabels"`
	CommonAnnotations KV `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

func GetPayloadData(reqBody io.ReadCloser) (DataPayload, error) {
	var payload DataPayload

	body, err := io.ReadAll(reqBody)
	if err != nil {
		fErr := fmt.Errorf("Could not read request body: %v", err)
		return payload, fErr
	}
	defer reqBody.Close()

	err = json.Unmarshal(body, &payload)
	if err != nil {
		fErr := fmt.Errorf("Could not parse request body: %v", err)
		return payload, fErr
	}

	return payload, nil
}
