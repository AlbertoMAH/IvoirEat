package main

// Fruit représente la structure de nos données de fruits dans la base de données.
type Fruit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
