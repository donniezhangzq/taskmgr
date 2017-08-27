package main

import (
    "database/sql"
    "flag"
    "fmt"
    "errors"

    "TaskDb"
    "FastSort"

    _ "github.com/mattn/go-sqlite3"
)


type t_type struct {
    id int
    type_name string
}

type t_item struct {
    id int
    item_name string
    status string
    pri int
    type_id int
}

func HandleType(
    TypeAddFlag * bool,
    TypeDelFlag *bool,
    TypeShowFlag *bool,
    TypeName *string,
    db *sql.DB) error {

    var err error
    if (!(*TypeAddFlag || *TypeDelFlag || *TypeShowFlag)) {
        help()
    } else if *TypeAddFlag && (!(*TypeDelFlag || *TypeShowFlag)){
    //type add
        if *TypeName == "" {
            return errors.New("未输入typename，请指定-t参数")
        }
        err = TaskDb.TypeAdd(*TypeName,db)
        return err
    } else if *TypeDelFlag && (!(*TypeAddFlag || *TypeShowFlag)) {
      // type del
        if *TypeName == "" {
            return errors.New("未输入typename，请指定-t参数")
        }
        err = TaskDb.TypeDel(*TypeName,db)
        return err
    } else if *TypeShowFlag && (!(*TypeAddFlag || *TypeDelFlag)) {
      // type show
        rows,err := TaskDb.TypeSelect(db)
        defer rows.Close()
        if err == nil {
            fmt.Println("type_id\ttype_name")
            for rows.Next() {
                var Ttype t_type
                err = rows.Scan(&Ttype.id, &Ttype.type_name)
                if err == nil {
                    fmt.Println(Ttype.id, Ttype.type_name)
                } else {
                    return err
                }
            }
        } else {
            return err
        }
    } else {
        return errors.New("show add del 只能同时执行一个")
    }
    return nil
}

func HandleItem(
    ItemAddFlag *bool,
    ItemGetFlag *bool,
    ItemDelFlag *bool,
    ItemUpdateFlag *bool,
    TypeName *string,
    ItemName *string,
    ItemStatus *string,
    ItemPri *int,
    db *sql.DB) error{

    var err error
    if *ItemAddFlag && !(*ItemGetFlag||*ItemDelFlag||*ItemUpdateFlag) {
        if *TypeName != "" && *ItemName != "" && *ItemPri !=0 {
           err = TaskDb.ItemAdd(*ItemName, *ItemPri, *TypeName, db)
           return err
        } else {
            return errors.New("缺少-t -i -p 参数")
        }
    } else if *ItemDelFlag && !(*ItemAddFlag||*ItemGetFlag||*ItemUpdateFlag) {
        if *TypeName != "" && *ItemName != "" {
            TaskDb.ItemDel(*ItemName, *TypeName, db)
        } else {
            return errors.New("缺少-t -i 参数")
        }
    } else if *ItemGetFlag && !(*ItemAddFlag||*ItemUpdateFlag||*ItemDelFlag) {
        var t []FastSort.T_item
        if *TypeName == "" {
            return errors.New("缺少-t参数")
        }
        TypeId,err := TaskDb.TypeGetId(*TypeName, db)
        if err != nil {
            return err
        }
        rows,err := TaskDb.ItemSelect(db)
        if err != nil {
            return err
        }
        defer rows.Close()
        for rows.Next() {
            var Titem FastSort.T_item
            err = rows.Scan(&Titem.Id,&Titem.Item_name,
                &Titem.Status,&Titem.Pri,&Titem.Type_id)
            if err != nil {
                return err
            } else {
                if Titem.Type_id == TypeId {
                    t = append(t, Titem)
                }
            }
        }
        item,err := FastSort.GetMaxPri(t)
        if  err != nil {
            fmt.Println(err)
            return err
        } else {
            fmt.Println("id\titem_name\tstatus\tpri")
            var output string
            output = fmt.Sprintf("%d\t%s\t%s\t%d", item.Id, item.Item_name, item.Status, item.Pri)
            fmt.Println(output)
        }
    } else if *ItemUpdateFlag && !(*ItemGetFlag||*ItemAddFlag||*ItemDelFlag) {
        if *TypeName != "" && *ItemName != "" {
           err = TaskDb.ItemUpdate(*ItemName, *TypeName, *ItemPri, *ItemStatus, db)
           return err
        } else {
            return errors.New("缺少-t -i 参数")
        }
    }
    return nil
}

func help() {
    var HelpStr string
    HelpStr = "任务管理工具task-go：工具理念：一次专注一件事，高效产出，减少打扰\n" +
        "-----------------------\n" +
        "----使用帮助：---------\n" +
        "tadd 增加类别，-tadd -t <type>\n" +
        "tdel 删除类别, -tdel -t <type>\n" +
        "tshow 展示类别, -tshow\n" +
        "---------------------------------------\n" +
        "iget 获取一个工作事项，-iget -t <type>\n" +
        "iupdate 设置一个工作事项的优先级，-iupdate -i <item> -t <type> (-p <pri>(0~5)) (-s <status>(todo/done))\n" +
        "iadd 增加事项，-iadd -i <item> -t <type> -p <pri>(default 1)\n" +
        "idel 删除事项，-idel -i <item> -t <type>\n" 
    fmt.Println(HelpStr)
}

func main() {
    var err error
    //connect db
    db, err := sql.Open("sqlite3", "./db/task.db")
    if err != nil {
        fmt.Println(err)
    }
    defer db.Close()
    //init flag
    var TypeAddFlag = flag.Bool("tadd", false, "add a type; usage: -type -add -t xxx")
    var TypeDelFlag = flag.Bool("tdel", false, "delete a type; usage: -type -del -t xxx")
    var TypeShowFlag = flag.Bool("tshow", false, "show all type names; usage: -type -show")
    var TypeName = flag.String("t", "", "type name")
    var ItemAddFlag = flag.Bool("iadd", false, "add a item")
    var ItemGetFlag = flag.Bool("iget", false, "get a most prioritized for a type")
    var ItemDelFlag = flag.Bool("idel", false, "delete item")
    var ItemUpdateFlag = flag.Bool("iupdate", false, "update item priority or status;" +
        " usage: -item -update (-s xxx)(-p xxx)")
    var ItemName = flag.String("i", "", "item name")
    var ItemStatus = flag.String("s", "", "item status: todo/done")
    var ItemPri = flag.Int("p", 0, "item priority")
    flag.Parse()

    //handle flag
    TypeFlag := (*TypeAddFlag || *TypeDelFlag || *TypeShowFlag)
    ItemFlag := (*ItemAddFlag || *ItemGetFlag || *ItemDelFlag || *ItemUpdateFlag)
    if (TypeFlag && ItemFlag) {
        fmt.Println("option error: only one option should be named")
    } else if  (!(TypeFlag || ItemFlag)) {
        help()
    } else if TypeFlag  {
        err = HandleType(TypeAddFlag, TypeDelFlag, TypeShowFlag, TypeName, db)
    } else if ItemFlag {
        err = HandleItem(ItemAddFlag, ItemGetFlag, ItemDelFlag, ItemUpdateFlag,
            TypeName ,ItemName, ItemStatus, ItemPri, db)
    }
    if err != nil {
        fmt.Println(err)
    }
}
