package main

import (
  "context"
  "fmt"
  "time"

  "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
  dsn := "postgres://postgres:cg123456.@localhost:5432/persona?sslmode=disable"
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  pool, err := pgxpool.New(ctx, dsn)
  if err != nil { panic(err) }
  defer pool.Close()

  rows, err := pool.Query(ctx, `
SELECT slug, status, is_active, source_character_slug, target_character_slug, work_slug
FROM public.pm_relations
WHERE slug = 'lin-daiyu-jia-baoyu-dream-of-the-red-chamber'
`)
  if err != nil { panic(err) }
  defer rows.Close()

  found := false
  for rows.Next() {
    found = true
    var slug, status, sourceSlug, targetSlug, workSlug string
    var isActive bool
    if err := rows.Scan(&slug, &status, &isActive, &sourceSlug, &targetSlug, &workSlug); err != nil { panic(err) }
    fmt.Println("relation", slug, status, isActive, sourceSlug, targetSlug, workSlug)
  }
  if err := rows.Err(); err != nil { panic(err) }
  if !found {
    fmt.Println("relation not found")
  }

  var count int
  err = pool.QueryRow(ctx, `
SELECT count(*)
FROM public.pm_relations
WHERE is_active = TRUE AND status = 'published' AND (source_character_slug = 'lin-daiyu' OR target_character_slug = 'lin-daiyu')
`).Scan(&count)
  if err != nil { panic(err) }
  fmt.Println("matching published active relation count", count)

  var cslug, cstatus string
  var cactive bool
  err = pool.QueryRow(ctx, `SELECT slug, status, is_active FROM public.pm_characters WHERE slug = 'lin-daiyu'`).Scan(&cslug, &cstatus, &cactive)
  if err != nil { panic(err) }
  fmt.Println("character", cslug, cstatus, cactive)
}
