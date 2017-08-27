package TaskDb

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "fmt"
    "errors"
)

func TypeAdd(type_name string, db *sql.DB) error {
    // insert type 
    var name string
    SqlStr := "select * from t_type where type_name=\"" + type_name + "\""
    err := db.QueryRow(SqlStr).Scan(&name)
    if err == nil{
        return errors.New("type_name"+type_name+"already exist")
    } else if err == sql.ErrNoRows {
        SqlStr := fmt.Sprintf("insert into t_type(type_name) values ('%s')", type_name)
        _, err := db.Exec(SqlStr)
        if err != nil {
            fmt.Println("insert err", err)
        }
    } else {
        return err
    }
    return nil
}

func TypeSelect(db *sql.DB) (*sql.Rows, error) {
    //query type
    SqlStr := "select * from t_type"
    result, err := db.Query(SqlStr)
    if err != nil {
        fmt.Println("seleect t_type error with:", err)
    }
    return result,err
}

func TypeDel(type_name string, db *sql.DB) error {
    // del a type
    SqlStr := fmt.Sprintf("delete from t_type where type_name='%s'", type_name)
    _, err := db.Exec(SqlStr)
    return err
}

func ItemAdd(item_name string, pri int, type_name string, db *sql.DB) error {
    // add a item;default status is todo
    status := "todo"
    type_id, err := TypeGetId(type_name, db)
    if err != nil {
        return err
    } else {
        err = ItemSelectExist(item_name, type_name, db)
        if err == sql.ErrNoRows {
            SqlStr := fmt.Sprintf("insert into t_item(item_name,status,pri,type_id)" +
                " values('%s','%s',%d,%d)", item_name, status, pri, type_id)
            _, err = db.Exec(SqlStr)
            return err
        } else if err != nil {
            return err
        } else {
            return err
        }
    }
}

func ItemSelectExist(item_name string, type_name string, db *sql.DB) error {
    var name string
    var type_id int
    type_id,err := TypeGetId(type_name, db)
    if err != nil {
        return err
    }
    SqlStr := fmt.Sprintf("select item_name from t_item where item_name='%s' and type_id=%d", item_name, type_id)
    err = db.QueryRow(SqlStr).Scan(&name)
    return err
}

func TypeGetId(type_name string, db *sql.DB) (int,error){
    var type_id int
    SqlStr := fmt.Sprintf("select type_id from t_type where type_name='%s'", type_name)
    err := db.QueryRow(SqlStr).Scan(&type_id)
    return type_id,err
}

func ItemDel(item_name string, type_name string, db *sql.DB) error {
    //delete a item
    type_id, err := TypeGetId(type_name, db)
    if err != nil {
        return err
    } else {
        SqlStr := fmt.Sprintf("delete from t_item where item_name='%s' and type_id=%d",
            item_name, type_id)
        _, err = db.Exec(SqlStr)
        return err
    }
}

func ItemSelect(db *sql.DB) (*sql.Rows, error) {
    //select all item
    SqlStr := fmt.Sprintf("select * from t_item")
    return db.Query(SqlStr)
}

func ItemUpdate(item_name string, type_name string, pri int, status string, db *sql.DB) error {
    // update status or pri
    //if pri is not 0,update pri; if status is not "",update status
    //status:todo/done
    err := ItemSelectExist(item_name, type_name, db)
    if err != nil {
        return err
    }
    type_id,err := TypeGetId(type_name, db)
    if err != nil {
        return err
    }
    if pri == 0 && status == "" {
        return errors.New("no need to update")
    } else {
        if pri != 0 {
            SqlStr := fmt.Sprintf("update t_item set pri=%d where item_name='%s' and type_id=%d", pri, item_name, type_id)
            _, err = db.Exec(SqlStr)
            if err != nil {
                return err
            }
        }
        if status != "" {
            SqlStr := fmt.Sprintf("update t_item set status='%s' where item_name='%s' and type_id=%d", status, item_name, type_id)
            _, err = db.Exec(SqlStr)
            if err != nil {
                return err
            }
        }
    }
    return nil
}

