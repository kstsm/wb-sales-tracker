package queries

const (
	CreateItemQuery = `
		INSERT INTO items (id,
		                   type,
		                   amount,
		                   date,
		                   category,
		                   created_at,
		                   updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
`

	GetItemByIDQuery = `
		SELECT id,
		       type,
		       amount,
		       date,
		       category,
		       created_at,
		       updated_at
		FROM items
		WHERE id = $1
`

	UpdateItemQuery = `
		UPDATE items
		SET 
			type = COALESCE($2, type),
			amount = COALESCE($3, amount),
			date = COALESCE($4, date),
			category = COALESCE($5, category),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id, type, amount, date, category, created_at, updated_at
`

	DeleteItemQuery = `
		DELETE FROM items
		WHERE id = $1
		RETURNING id
`

	BaseCountQuery = `
    SELECT COUNT(*)
    FROM items
`

	BaseSelectQuery = `
    SELECT id,
           type,
           amount,
           date,
           category,
           created_at,
           updated_at
    FROM items
`

	AnalyticsQuery = `
		SELECT 
			COALESCE(SUM(amount), 0) as sum,
			AVG(amount) as avg,
			COUNT(*) as count,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) as median,
			PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount) as percentile_90
		FROM items
		%s
	`

	AnalyticsGroupedByDayQuery = `
		SELECT 
			DATE(date) as group_key,
			COALESCE(SUM(amount), 0) as sum,
			AVG(amount) as avg,
			COUNT(*) as count,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) as median,
			PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount) as percentile_90
		FROM items
		%s
		GROUP BY DATE(date)
		ORDER BY DATE(date)
	`

	AnalyticsGroupedByWeekQuery = `
		SELECT 
			DATE_TRUNC('week', date) as group_key,
			COALESCE(SUM(amount), 0) as sum,
			AVG(amount) as avg,
			COUNT(*) as count,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) as median,
			PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount) as percentile_90
		FROM items
		%s
		GROUP BY DATE_TRUNC('week', date)
		ORDER BY DATE_TRUNC('week', date)
	`

	AnalyticsGroupedByCategoryQuery = `
		SELECT 
			category as group_key,
			COALESCE(SUM(amount), 0) as sum,
			AVG(amount) as avg,
			COUNT(*) as count,
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY amount) as median,
			PERCENTILE_CONT(0.9) WITHIN GROUP (ORDER BY amount) as percentile_90
		FROM items
		%s
		GROUP BY category
		ORDER BY category
	`
)
