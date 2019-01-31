package main

// language=GraphQL
var graphQLSchema = `
schema {
	query: Query
}

type Query {
	search(query: String): [ImageGroup!]!
}

type ImageGroup {
	id: ID // its in fact Sha256 taken from image that has proper semver assigned
	name: String!
	version: String!
	images: [Image!]!
	checksum: Checksum!
	createdAt: String
	dockerVersion: String
	operatingSystem: String
	author: String
	preRelease: String
}

type Image {
    id: ID
	name: String!
	version: String!
	manifestUrl: String!
    path: String!
	pullUrl: String!
	author: String
	createdAt: String
}

type Checksum {
	sha1: String!
	sha256: String!
}
`
