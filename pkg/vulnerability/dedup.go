package vulnerability

type Deduplicator struct {
    seen map[string]bool
}

func NewDeduplicator() *Deduplicator {
    return &Deduplicator{
        seen: make(map[string]bool),
    }
}

func (d *Deduplicator) Deduplicate(vulns []Vulnerability) []Vulnerability {
    var result []Vulnerability
    key := func(v Vulnerability) string {
        return v.URL + "|" + v.TemplateID
    }
    
    for _, v := range vulns {
        k := key(v)
        if !d.seen[k] {
            d.seen[k] = true
            result = append(result, v)
        }
    }
    return result
}

func (d *Deduplicator) Clear() {
    d.seen = make(map[string]bool)
}

func (d *Deduplicator) Count() int {
    return len(d.seen)
}
