node {
    jobParts = env.JOB_NAME.split('/')
    golang(
            project: jobParts[0].toLowerCase(),
            name: jobParts[1].toLowerCase(),
            library: true
    )
}
