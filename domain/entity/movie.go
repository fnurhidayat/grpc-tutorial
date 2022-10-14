package entity

type Movie struct {
	Id      uint32 `db:"id"`
	Title   string `db:"title"`
	Summary string `db:"summary"`
	Rating  uint32 `db:"rating"`
}
