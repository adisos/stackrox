export function mockListComplianceProfiles() {
    return [
        {
            id: '1',
            name: 'CIS Docker',
            profile_version: '1.0',
            product_type: ['Container'],
            standard: 'CIS',
            description: 'Docker containers based on CIS benchmarks.',
            rules: [],
            product: 'Docker',
            title: 'CIS Benchmark for Docker',
        },
        {
            id: '2',
            name: 'CIS K8s',
            profile_version: '1.0',
            product_type: ['Container Orchestration'],
            standard: 'CIS',
            description: 'Kubernetes clusters following CIS benchmarks.',
            rules: [],
            product: 'Kubernetes',
            title: 'CIS Benchmark for Kubernetes',
        },
        {
            id: '3',
            name: 'HIPPA',
            profile_version: '1.0',
            product_type: ['Healthcare'],
            standard: 'HIPAA',
            description: 'healthcare applications ensuring HIPAA compliance.',
            rules: [],
            product: 'Healthcare Systems',
            title: 'HIPAA Compliance Profile',
        },
        {
            id: '4',
            name: 'NIST SP 800-190',
            profile_version: '1.0',
            product_type: ['Cloud Computing', 'Container Security'],
            standard: 'NIST',
            description:
                'security in cloud computing and container technologies as outlined in NIST SP 800-190.',
            rules: [],
            product: 'Cloud Services',
            title: 'NIST SP 800-190 Compliance',
        },
        {
            id: '5',
            name: 'PCI',
            profile_version: '1.0',
            product_type: ['Payment Processing'],
            standard: 'PCI DSS',
            description: 'payment card industry data security standard (PCI DSS) compliance.',
            rules: [],
            product: 'Payment Systems',
            title: 'PCI DSS Compliance Profile',
        },
    ];
}