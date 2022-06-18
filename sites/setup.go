package main

import (
  "strings"
  "path/filepath"
  // "github.com/saenuma/flaarum"
  "github.com/saenuma/flaarum/flaarum_shared"
  "github.com/bankole7782/office683/office683_shared"
)


func findIn(container []string, toFind string) int {
  for i, inContainer := range container {
    if inContainer == toFind {
      return i
    }
  }
  return -1
}


func formatTableStruct(tableStruct flaarum_shared.TableStruct) string {
	stmt := "table: " + tableStruct.TableName + "\n"
	stmt += "fields:\n"
	for _, fieldStruct := range tableStruct.Fields {
		stmt += "  " + fieldStruct.FieldName + " " + fieldStruct.FieldType
		if fieldStruct.Required {
			stmt += " required"
		}
		if fieldStruct.Unique {
			stmt += " unique"
		}
		if fieldStruct.NotIndexed {
			stmt += " nindex"
		}
		stmt += "\n"
	}
	stmt += "::\n"
	if len(tableStruct.ForeignKeys) > 0 {
		stmt += "foreign_keys:\n"
		for _, fks := range tableStruct.ForeignKeys {
			stmt += "  " + fks.FieldName + " " + fks.PointedTable + " " + fks.OnDelete + "\n"
		}
		stmt += "::\n"
	}

	if len(tableStruct.UniqueGroups) > 0 {
		stmt += "unique_groups:\n"
		for _, ug := range tableStruct.UniqueGroups {
			stmt += "  " + strings.Join(ug, " ") + "\n"
		}
		stmt += "::\n"
	}

	return stmt
}


func createOrUpdateTable(stmt string) error {
  FRCL := office683_shared.GetFlaarumClient()

	tables, err := FRCL.ListTables()
	if err != nil {
		return err
	}

	tableStruct, err := flaarum_shared.ParseTableStructureStmt(stmt)
	if err != nil {
		return err
	}
	if findIn(tables, tableStruct.TableName) == -1 {
		// table doesn't exist
		err = FRCL.CreateTable(stmt)
		if err != nil {
			return err
		}
	} else {
		// table exists check if it needs update
    currentVersionNum, err := FRCL.GetCurrentTableVersionNum(tableStruct.TableName)
    if err != nil {
      return err
    }

		oldStmt, err := FRCL.GetTableStructure(tableStruct.TableName, currentVersionNum)
		if err != nil {
			return err
		}

		if oldStmt != formatTableStruct(tableStruct) {
			err = FRCL.UpdateTableStructure(stmt)
			if err != nil {
				return err
			}
		}

	}
	return nil
}


func CreateOrUpdateAllTables() error {
  dirFIs, err := office683_shared.FlaarumStmts.ReadDir("flaarum_stmts")
  if err != nil {
    return err
  }

  orderOfTables := []string{
    "users", "teams", "team_members", "sessions", "events", "docs_folders", "docs_images",
    "docs", "cab_folders", "cab_files"
  }
  for _, tableName := range orderOfTables {
    stmt, err := office683_shared.FlaarumStmts.ReadFile(filepath.Join("flaarum_stmts", tableName + ".txt"))
    if err != nil {
      return err
    }

    err = createOrUpdateTable(string(stmt))
    if err != nil {
      return err
    }
  }

  return nil
}
