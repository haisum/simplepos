package items

import (
	"github.com/haisum/simplepos/db"
	"github.com/pkg/errors"
	"gopkg.in/doug-martin/goqu.v4"
)

//Item represents single Item which is a good that's sold or bought
type Item struct {
	ID         int64  `db:"ID" goqu:"skipinsert"`
	Name       string `db:"Name"`
	Attributes string `db:"Attributes"`
	SellPrice  uint64 `db:"SellPrice"`
	BuyPrice   uint64 `db:"BuyPrice"`
	Stock      int64  `db:"Stock"`
	Comments   string `db:"Comments"`
}

//Validate validates an item
func (t *Item) Validate() error {
	if t.Name == "" {
		return errors.New("Name is a required field.")
	}
	return nil
}

//Criteria specifies conditions to filter items
type Criteria struct {
	ID               int64
	Name             string
	StockLessThan    int64
	StockGreaterThan int64
}

// List filters out items
func List(c Criteria) ([]Item, error) {
	list := []Item{}
	ds := db.Get().From("Item")
	if c.ID != 0 {
		ds = ds.Where(goqu.I("ID").Eq(c.ID))
	}
	if c.Name != "" {
		ds = ds.Where(goqu.I("Name").Like("%" + c.Name + "%"))
	}
	if c.StockLessThan != 0 {
		ds = ds.Where(goqu.I("Stock").Lt(c.StockLessThan))
	}
	if c.StockGreaterThan != 0 {
		ds = ds.Where(goqu.I("Stock").Gt(c.StockGreaterThan))
	}
	query, _, _ := ds.ToSql()
	err := db.Get().ScanStructs(&list, query)
	return list, errors.Wrap(err, "Error in quering db for items.")
}

// Add adds items, records shouldn't specify ID when inserting. Last affected and error is returned.
func Add(items []Item) (int64, error) {
	for _, it := range items {
		if err := it.Validate(); err != nil {
			return 0, errors.Wrapf(err, "Error in insert. %+v", it)
		}
	}
	rs, err := db.Get().From("Item").Insert(items).Exec()
	if err != nil {
		return 0, errors.Wrap(err, "Error in inserting records")
	}
	affected, _ := rs.RowsAffected()
	return affected, errors.Wrap(err, "Error in inserting records")
}

// Update updates items, each Item must have ID field when updating. Last affected and error is returned.
func Update(items []Item) (int64, error) {
	for _, it := range items {
		if err := it.Validate(); err != nil {
			return 0, errors.Wrapf(err, "Error in update. %+v", it)
		}
		if it.ID <= 0 {
			return 0, errors.New("ID is required when updating a record.")
		}
	}
	totalAffected := int64(0)
	tx, err := db.Get().Begin()
	if err != nil {
		return 0, errors.Wrap(err, "Error in beginning transaction")
	}
	for _, it := range items {
		rs, err := tx.From("Item").Where(goqu.I("ID").Eq(it.ID)).Update(it).Exec()
		if err != nil {
			if rxErr := tx.Rollback(); rxErr != nil {
				return totalAffected, errors.Wrap(rxErr, "Error in rollback")
			}
			return 0, errors.Wrap(err, "Error in updating. Rolled back.")
		}
		affected, err := rs.RowsAffected()
		if err != nil {
			return 0, errors.Wrap(err, "Error in getting affected rows")
		}
		totalAffected = affected + totalAffected
	}
	err = tx.Commit()
	return totalAffected, errors.Wrap(err, "Error in commit.")
}

//Delete deletes items with given primary keys. Last affected and error is returned.
func Delete(ids []int64) (int64, error) {
	tx, err := db.Get().Begin()
	if err != nil {
		return 0, errors.Wrap(err, "Error in beginning transaction.")
	}
	rs, err := tx.From("Item").Where(goqu.I("ID").In(ids)).Delete().Exec()
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "Error in delete, rolling back.")
	}
	tx.Commit()
	return rs.RowsAffected()
}
