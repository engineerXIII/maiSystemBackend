package repository

const (
	createProduct = `INSERT INTO products (product_name, color, description, factory, cost, created_at)
						VALUES ($1, $2, $3, $4, $5, now())
						RETURNING *`
	updateProduct = `UPDATE products
						SET product_name = COALESCE(NULLIF($1, ''), product_name),
							color = COALESCE(NULLIF($2, ''), color),
							description = COALESCE(NULLIF($3, ''), description),
							factory = COALESCE(NULLIF($4, ''), factory),
							cost = COALESCE(NULLIF($5, ''), cost),
							updated_at = now()
						WHERE product_id = $6
						RETURNING *`
	getProductByID = `SELECT product_id, 
						product_name,
						color,
						description,
						factory,
						cost,
						updated_at
					FROM products
					WHERE product_id = $1`
	deleteProduct = `DELETE FROM products WHERE product_id = $1`
	getTotalCount = `SELECT COUNT(product_id) FROM products`
	getProducts   = `SELECT product_id, 
						product_name,
						color,
						description,
						factory,
						cost,
						created_at,
						updated_at
					FROM products
					ORDER BY created_at, updated_at OFFSET $1 LIMIT $2`
	findByNameCount = `SELECT COUNT(*)
					FROM products
					WHERE product_name ILIKE '%' || $1 || '%'`
	findByName = `SELECT product_id, 
						product_name,
						color,
						description,
						factory,
						cost,
						created_at,
						updated_at
					FROM products
					WHERE product_name ILIKE '%' || $1 || '%'
					ORDER BY product_name, created_at, updated_at
					OFFSET $2 LIMIT $3`
)
