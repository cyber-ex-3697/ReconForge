package takeover

type ServiceFingerprint struct {
    Service   string
    Fingerprint string
    Vulnerable bool
}

type FingerprintMatcher struct {
    fingerprints []ServiceFingerprint
}

func NewFingerprintMatcher() *FingerprintMatcher {
    return &FingerprintMatcher{
        fingerprints: []ServiceFingerprint{
            {Service: "AWS S3", Fingerprint: "NoSuchBucket", Vulnerable: true},
            {Service: "GitHub Pages", Fingerprint: "There isn't a GitHub Pages site here", Vulnerable: true},
            {Service: "Heroku", Fingerprint: "No such app", Vulnerable: true},
            {Service: "Azure Blob", Fingerprint: "The specified blob does not exist", Vulnerable: true},
            {Service: "CloudFront", Fingerprint: "The request could not be satisfied", Vulnerable: true},
            {Service: "Readme.io", Fingerprint: "Project doesnt exist", Vulnerable: true},
            {Service: "Surge.sh", Fingerprint: "project not found", Vulnerable: true},
            {Service: "Bitbucket", Fingerprint: "Repository not found", Vulnerable: true},
            {Service: "Pantheon", Fingerprint: "No site found", Vulnerable: true},
            {Service: "Fastly", Fingerprint: "Fastly error: unknown domain", Vulnerable: true},
        },
    }
}

func (f *FingerprintMatcher) Match(response string) *ServiceFingerprint {
    for _, fp := range f.fingerprints {
        if f.contains(response, fp.Fingerprint) {
            return &fp
        }
    }
    return nil
}

func (f *FingerprintMatcher) contains(response, pattern string) bool {
    return len(response) > 0 && len(pattern) > 0 && 
           (response == pattern || 
            (len(response) > len(pattern) && 
             response[len(response)-len(pattern):] == pattern))
}

func (f *FingerprintMatcher) AddFingerprint(service, fingerprint string, vulnerable bool) {
    f.fingerprints = append(f.fingerprints, ServiceFingerprint{
        Service:     service,
        Fingerprint: fingerprint,
        Vulnerable:  vulnerable,
    })
}

func (f *FingerprintMatcher) IsVulnerable(response string) bool {
    match := f.Match(response)
    return match != nil && match.Vulnerable
}
