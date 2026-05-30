package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"

	"ds-backend/internal/domain/models"
	"ds-backend/internal/domain/ports"
)

type PropertyRepository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewPropertyRepository(db *sql.DB) ports.PropertyRepository {
	return &PropertyRepository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}
}

func (r *PropertyRepository) GetAll(ctx context.Context) ([]models.PropertyCard, error) {
	query := r.qb.Select(
		"p.id", "p.title", "p.price", "p.bedrooms", "p.bathrooms",
		"p.construction_m2", "p.city", "p.neighborhood",
		"p.property_type", "p.transaction_type", "p.status", "i.path", "p.created_at",
		"u_owner.id", "u_owner.username", "u_owner.email", "u_owner.notes",
	).
		From("properties p").
		LeftJoin("images i ON i.property_id = p.id AND i.main_image = 1 AND i.deleted_at IS NULL").
		LeftJoin("users u_owner ON p.owner_id = u_owner.id AND u_owner.deleted_at IS NULL").
		Where(squirrel.Expr("p.deleted_at IS NULL"))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting all properties")
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for getting all properties")
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close rows after getting all properties")
		}
	}()

	var properties []models.PropertyCard
	var propertyIds []uint
	propertyMap := make(map[uint]models.PropertyCard)

	for rows.Next() {
		var property models.PropertyCard
		var ownerID sql.NullInt32
		var ownerUsername sql.NullString
		var ownerEmail sql.NullString
		var ownerNotes sql.NullString

		if err := rows.Scan(
			&property.ID,
			&property.Title,
			&property.Price,
			&property.Bedrooms,
			&property.Bathrooms,
			&property.ConstructionM2,
			&property.City,
			&property.Neighborhood,
			&property.PropertyType,
			&property.TransactionType,
			&property.Status,
			&property.MainImagePath,
			&property.CreatedAt,
			&ownerID,
			&ownerUsername,
			&ownerEmail,
			&ownerNotes,
		); err != nil {
			logrus.WithError(err).Error("Failed to scan property row")
			return nil, err
		}

		// Set owner information if available
		if ownerID.Valid {
			property.OwnerID = &models.UserInfo{
				ID:       uint(ownerID.Int32),
				Username: ownerUsername.String,
				Email:    ownerEmail.String,
				Notes:    (*string)(nil),
			}
			if ownerNotes.Valid {
				property.OwnerID.Notes = &ownerNotes.String
			}
		}

		// Initialize agent as empty slice
		property.Agent = make([]models.UserInfo, 0)

		propertyIds = append(propertyIds, property.ID)
		propertyMap[property.ID] = property
	}
	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("Error occurred while iterating over property rows")
		return nil, err
	}

	if len(propertyMap) == 0 {
		logrus.Warn("No properties found in the database")
		return []models.PropertyCard{}, nil
	}

	// Query all agents for these properties
	agentQuery := r.qb.Select(
		"up.property_id", "u.id", "u.username", "u.email", "u.notes",
	).
		From("user_properties up").
		Join("users u ON up.user_id = u.id AND u.deleted_at IS NULL").
		Where(squirrel.Eq{"up.property_id": propertyIds})

	agentSqlStr, agentArgs, err := agentQuery.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting property agents")
		return nil, err
	}

	agentRows, err := r.db.QueryContext(ctx, agentSqlStr, agentArgs...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for getting property agents")
		return nil, err
	}
	defer func() {
		if err := agentRows.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close rows after getting property agents")
		}
	}()

	for agentRows.Next() {
		var propID uint
		var agent models.UserInfo
		var agentNotes sql.NullString

		if err := agentRows.Scan(&propID, &agent.ID, &agent.Username, &agent.Email, &agentNotes); err != nil {
			logrus.WithError(err).Error("Failed to scan agent row")
			return nil, err
		}

		if agentNotes.Valid {
			agent.Notes = &agentNotes.String
		}

		if prop, exists := propertyMap[propID]; exists {
			prop.Agent = append(prop.Agent, agent)
			propertyMap[propID] = prop
		}
	}

	if err := agentRows.Err(); err != nil {
		logrus.WithError(err).Error("Error occurred while iterating over agent rows")
		return nil, err
	}

	// Reconstruct properties array from map in consistent order
	properties = make([]models.PropertyCard, 0, len(propertyMap))
	for _, propID := range propertyIds {
		if prop, exists := propertyMap[propID]; exists {
			properties = append(properties, prop)
		}
	}

	return properties, nil
}

func (r *PropertyRepository) GetPropertyCardByID(ctx context.Context, id uint) (*models.PropertyCard, error) {
	query := r.qb.Select(
		"p.id", "p.title", "p.price", "p.bedrooms", "p.bathrooms",
		"p.construction_m2", "p.city", "p.neighborhood",
		"p.property_type", "p.transaction_type", "p.status", "i.path", "p.created_at",
		"u_owner.id", "u_owner.username", "u_owner.email", "u_owner.notes",
	).
		From("properties p").
		LeftJoin("images i ON i.property_id = p.id AND i.main_image = 1 AND i.deleted_at IS NULL").
		LeftJoin("users u_owner ON p.owner_id = u_owner.id AND u_owner.deleted_at IS NULL").
		Where(squirrel.Eq{"p.id": id}).
		Where(squirrel.Expr("p.deleted_at IS NULL"))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting property card by ID")
		return nil, err
	}

	var property models.PropertyCard
	var ownerID sql.NullInt32
	var ownerUsername sql.NullString
	var ownerEmail sql.NullString
	var ownerNotes sql.NullString

	err = r.db.QueryRowContext(ctx, sqlStr, args...).Scan(
		&property.ID,
		&property.Title,
		&property.Price,
		&property.Bedrooms,
		&property.Bathrooms,
		&property.ConstructionM2,
		&property.City,
		&property.Neighborhood,
		&property.PropertyType,
		&property.TransactionType,
		&property.Status,
		&property.MainImagePath,
		&property.CreatedAt,
		&ownerID,
		&ownerUsername,
		&ownerEmail,
		&ownerNotes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithError(err).Warnf("No property found with ID %d", id)
			return nil, nil
		}
		logrus.WithError(err).Error("Failed to execute query for getting property card by ID")
		return nil, err
	}

	// Set owner information if available
	if ownerID.Valid {
		property.OwnerID = &models.UserInfo{
			ID:       uint(ownerID.Int32),
			Username: ownerUsername.String,
			Email:    ownerEmail.String,
			Notes:    (*string)(nil),
		}
		if ownerNotes.Valid {
			property.OwnerID.Notes = &ownerNotes.String
		}
	}

	agentQuery := r.qb.Select(
		"u.id", "u.username", "u.email", "u.notes",
	).
		From("user_properties up").
		Join("users u ON up.user_id = u.id AND u.deleted_at IS NULL").
		Where(squirrel.Eq{"up.property_id": id})

	agentSqlStr, agentArgs, err := agentQuery.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting property agents")
		return nil, err
	}

	agentRows, err := r.db.QueryContext(ctx, agentSqlStr, agentArgs...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for getting property agents")
		return nil, err
	}
	defer func() {
		if err := agentRows.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close rows after getting property agents")
		}
	}()

	// Initialize agents as empty slice to ensure consistent JSON output
	agents := make([]models.UserInfo, 0)
	for agentRows.Next() {
		var agent models.UserInfo
		var agentNotes sql.NullString

		if err := agentRows.Scan(&agent.ID, &agent.Username, &agent.Email, &agentNotes); err != nil {
			logrus.WithError(err).Error("Failed to scan agent row")
			return nil, err
		}

		if agentNotes.Valid {
			agent.Notes = &agentNotes.String
		}

		agents = append(agents, agent)
	}

	if err := agentRows.Err(); err != nil {
		logrus.WithError(err).Error("Error occurred while iterating over agent rows")
		return nil, err
	}

	property.Agent = agents

	return &property, nil
}

func (r *PropertyRepository) GetByID(ctx context.Context, id uint) (*models.PropertyResponse, error) {
	// Query property with owner information
	query := r.qb.Select(
		"p.id", "p.title", "p.address", "p.neighborhood", "p.city", "p.zone",
		"p.reference", "p.price", "p.construction_m2", "p.land_m2", "p.is_occupied",
		"p.is_furnished", "p.floors", "p.bedrooms", "p.bathrooms", "p.garage_size",
		"p.garden_m2", "p.gas_types", "p.amenities", "p.extras", "p.utilities",
		"p.notes", "p.description", "p.owner_id", "p.property_type",
		"p.transaction_type", "p.status", "p.created_at", "p.updated_at",
		"u_owner.id", "u_owner.username", "u_owner.email", "u_owner.notes",
	).
		From("properties p").
		LeftJoin("users u_owner ON p.owner_id = u_owner.id AND u_owner.deleted_at IS NULL").
		Where(squirrel.And{
			squirrel.Eq{"p.id": id},
			squirrel.Expr("p.deleted_at IS NULL"),
		})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting property by ID")
		return nil, err
	}

	var property models.Property
	var ownerID sql.NullInt32
	var ownerUsername sql.NullString
	var ownerEmail sql.NullString
	var ownerNotes sql.NullString

	err = r.db.QueryRowContext(ctx, sqlStr, args...).Scan(
		&property.ID,
		&property.Title,
		&property.Address,
		&property.Neighborhood,
		&property.City,
		&property.Zone,
		&property.Reference,
		&property.Price,
		&property.ConstructionM2,
		&property.LandM2,
		&property.IsOccupied,
		&property.IsFurnished,
		&property.Floors,
		&property.Bedrooms,
		&property.Bathrooms,
		&property.GarageSize,
		&property.GardenM2,
		&property.GasTypes,
		&property.Amenities,
		&property.Extras,
		&property.Utilities,
		&property.Notes,
		&property.Description,
		&property.OwnerID,
		&property.PropertyType,
		&property.TransactionType,
		&property.Status,
		&property.CreatedAt,
		&property.UpdatedAt,
		&ownerID,
		&ownerUsername,
		&ownerEmail,
		&ownerNotes,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithError(err).Warnf("No property found with ID %d", id)
			return nil, err
		}
		logrus.WithError(err).Error("Failed to execute query for getting property by ID")
		return nil, err
	}

	// Convert to response
	response := property.ToResponse()

	// Set owner information if available
	if ownerID.Valid {
		response.Owner = &models.UserInfo{
			ID:       uint(ownerID.Int32),
			Username: ownerUsername.String,
			Email:    ownerEmail.String,
			Notes:    (*string)(nil),
		}
		if ownerNotes.Valid {
			response.Owner.Notes = &ownerNotes.String
		}
	}

	agentQuery := r.qb.Select(
		"u.id", "u.username", "u.email", "u.notes",
	).
		From("user_properties up").
		Join("users u ON up.user_id = u.id AND u.deleted_at IS NULL").
		Where(squirrel.Eq{"up.property_id": id})

	agentSqlStr, agentArgs, err := agentQuery.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for getting property agents")
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, agentSqlStr, agentArgs...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for getting property agents")
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.WithError(err).Error("Failed to close rows after getting property agents")
		}
	}()

	var agents []models.UserInfo
	for rows.Next() {
		var agent models.UserInfo
		var agentNotes sql.NullString

		if err := rows.Scan(&agent.ID, &agent.Username, &agent.Email, &agentNotes); err != nil {
			logrus.WithError(err).Error("Failed to scan agent row")
			return nil, err
		}

		if agentNotes.Valid {
			agent.Notes = &agentNotes.String
		}

		agents = append(agents, agent)
	}

	if err := rows.Err(); err != nil {
		logrus.WithError(err).Error("Error occurred while iterating over agent rows")
		return nil, err
	}

	response.Agent = agents

	return response, nil
}

func (r *PropertyRepository) Create(ctx context.Context, property *models.Property) (*models.PropertyResponse, error) {
	// Start transaction for atomic operations
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to begin transaction for creating property")
		return nil, err
	}
	defer tx.Rollback()

	query := r.qb.Insert("properties").
		Columns(
			"title", "address", "neighborhood", "city",
			"zone", "reference", "price", "construction_m2", "land_m2",
			"is_occupied", "is_furnished", "floors", "bedrooms", "bathrooms",
			"garage_size", "garden_m2", "gas_types", "amenities", "extras",
			"utilities", "notes", "description", "owner_id", "property_type",
			"transaction_type", "status",
		).
		Values(
			property.Title, property.Address, property.Neighborhood, property.City,
			property.Zone, property.Reference, property.Price, property.ConstructionM2, property.LandM2,
			property.IsOccupied, property.IsFurnished, property.Floors, property.Bedrooms, property.Bathrooms,
			property.GarageSize, property.GardenM2, property.GasTypes, property.Amenities, property.Extras,
			property.Utilities, property.Notes, property.Description, property.OwnerID, property.PropertyType,
			property.TransactionType, property.Status,
		)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for creating a new property")
		return nil, err
	}

	result, err := tx.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for creating a new property")
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		logrus.WithError(err).Error("Failed to get last insert ID")
		return nil, err
	}

	if id < 0 || id > int64(^uint(0)>>1) {
		return nil, errors.New("invalid ID: integer overflow")
	}
	property.ID = uint(id)

	// Insert all agents into user_properties junction table
	for _, agentID := range property.UserID {
		agentQuery := r.qb.Insert("user_properties").
			Columns("user_id", "property_id").
			Values(agentID, property.ID)

		agentSqlStr, agentArgs, err := agentQuery.ToSql()
		if err != nil {
			logrus.WithError(err).Error("Failed to build SQL query for inserting agent to user_properties")
			return nil, err
		}

		_, err = tx.ExecContext(ctx, agentSqlStr, agentArgs...)
		if err != nil {
			logrus.WithError(err).Error("Failed to insert agent into user_properties")
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.WithError(err).Error("Failed to commit transaction for creating property")
		return nil, err
	}

	property.CreatedAt = time.Now()
	logrus.Infof("Property created successfully with ID: %d with %d agents", property.ID, len(property.UserID))
	return property.ToResponse(), nil
}

func (r *PropertyRepository) Update(ctx context.Context, property *models.Property) (*models.PropertyResponse, error) {
	// Start transaction for atomic operations
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to begin transaction for updating property")
		return nil, err
	}
	defer tx.Rollback()

	query := r.qb.Update("properties").
		Set("title", property.Title).
		Set("address", property.Address).
		Set("neighborhood", property.Neighborhood).
		Set("city", property.City).
		Set("zone", property.Zone).
		Set("reference", property.Reference).
		Set("price", property.Price).
		Set("construction_m2", property.ConstructionM2).
		Set("land_m2", property.LandM2).
		Set("is_occupied", property.IsOccupied).
		Set("is_furnished", property.IsFurnished).
		Set("floors", property.Floors).
		Set("bedrooms", property.Bedrooms).
		Set("bathrooms", property.Bathrooms).
		Set("garage_size", property.GarageSize).
		Set("garden_m2", property.GardenM2).
		Set("gas_types", property.GasTypes).
		Set("amenities", property.Amenities).
		Set("extras", property.Extras).
		Set("utilities", property.Utilities).
		Set("notes", property.Notes).
		Set("description", property.Description).
		Set("owner_id", property.OwnerID).
		Set("property_type", property.PropertyType).
		Set("transaction_type", property.TransactionType).
		Set("status", property.Status).
		Where(squirrel.Eq{"id": property.ID}).
		Where(squirrel.Expr("deleted_at IS NULL"))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for updating a property")
		return nil, err
	}

	result, err := tx.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for updating a property")
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Failed to get rows affected after updating a property")
		return nil, err
	}

	if rowsAffected == 0 {
		logrus.Warnf("No rows updated for Property ID %d. It may not exist or has already been deleted.", property.ID)
		return nil, errors.New("property not found or already deleted")
	}

	deleteQuery := r.qb.Delete("user_properties").
		Where(squirrel.Eq{"property_id": property.ID})

	deleteSqlStr, deleteArgs, err := deleteQuery.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for deleting old agents")
		return nil, err
	}

	_, err = tx.ExecContext(ctx, deleteSqlStr, deleteArgs...)
	if err != nil {
		logrus.WithError(err).Error("Failed to delete old agent relationships")
		return nil, err
	}

	// Insert new agent relationships
	for _, agentID := range property.UserID {
		agentQuery := r.qb.Insert("user_properties").
			Columns("user_id", "property_id").
			Values(agentID, property.ID)

		agentSqlStr, agentArgs, err := agentQuery.ToSql()
		if err != nil {
			logrus.WithError(err).Error("Failed to build SQL query for inserting agent to user_properties")
			return nil, err
		}

		_, err = tx.ExecContext(ctx, agentSqlStr, agentArgs...)
		if err != nil {
			logrus.WithError(err).Error("Failed to insert agent into user_properties")
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		logrus.WithError(err).Error("Failed to commit transaction for updating property")
		return nil, err
	}

	logrus.Infof("Property with ID %d updated successfully with %d agents", property.ID, len(property.UserID))
	return property.ToResponse(), nil
}

func (r *PropertyRepository) Delete(ctx context.Context, id uint) error {
	query := r.qb.Update("properties").
		Set("deleted_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.Expr("deleted_at IS NULL"))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logrus.WithError(err).Error("Failed to build SQL query for deleting a property")
		return err
	}

	result, err := r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		logrus.WithError(err).Error("Failed to execute query for deleting a property")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Failed to get rows affected after deleting a property")
		return err
	}

	if rowsAffected == 0 {
		logrus.Warnf("No property found with ID %d or already deleted", id)
		return errors.New("property not found or already deleted")
	}

	return nil
}
