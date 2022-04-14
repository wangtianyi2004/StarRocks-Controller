package checkStatus

import(
    "fmt"
    "strings"
    "strconv"
    "sr-controller/sr-utl"
    "sr-controller/module"
    "database/sql"
)

type BeStatusStruct struct{

    BackendId                int
    Cluster                  string
    IP                       string
    HeartbeatServicePort     int
    BePort                   int
    HttpPort                 int
    BrpcPort                 int
    LastStartTime            sql.NullString
    LastHeartbeat            sql.NullString
    Alive                    bool
    SystemDecommissioned     bool
    ClusterDecommissioned    bool
    TabletNum                int
    DataUsedCapacity         string
    AvailCapacity            string
    TotalCapacity            sql.NullString
    UsedPct                  string
    MaxDiskUsedPct           string
    ErrMsg                   sql.NullString
    Version                  sql.NullString
    Status                   sql.NullString
    DataTotalCapacity        sql.NullString
    DataUsedPct              sql.NullString

}


var GBeStatArr []BeStatusStruct

func CheckBePortStatus(beId int) (checkPortRes bool, err error) {

    var infoMess string

    tmpUser := module.GYamlConf.Global.User
    tmpKeyRsa := module.GSshKeyRsa
    tmpBeHost := module.GYamlConf.BeServers[beId].Host
    tmpSshPort := module.GYamlConf.BeServers[beId].SshPort
    tmpHeartbeatServicePort := module.GYamlConf.BeServers[beId].HeartbeatServicePort
    checkCMD := fmt.Sprintf("netstat -nltp | grep ':%d '", tmpHeartbeatServicePort)

    output, err := utl.SshRun(tmpUser, tmpKeyRsa, tmpBeHost, tmpSshPort, checkCMD)

    if err != nil {
        infoMess = fmt.Sprintf("Error in run cmd when check BE port status [BeHost = %s, error = %v]", tmpBeHost, err)
        utl.Log("DEBUG", infoMess)
        return false, err
    }

    if strings.Contains(string(output), ":" + strconv.Itoa(tmpHeartbeatServicePort)) {
        infoMess = fmt.Sprintf("Check the BE query port %s:%d run successfully", tmpBeHost, tmpHeartbeatServicePort)
        utl.Log("DEBUG", infoMess)
        return true, nil
    }

    return false, err
}


func GetBeStatJDBC(beId int) (beStat BeStatusStruct, err error) {

    var infoMess string
    var tmpBeStat BeStatusStruct
    //GJdbcUser = "root"
    //GJdbcPasswd = ""
    //GJdbcDb = ""
    queryCMD := "show backends"
    tmpBeHost := module.GYamlConf.BeServers[beId].Host
    tmpHeartbeatServicePort := module.GYamlConf.BeServers[beId].HeartbeatServicePort

    rows, err := utl.RunSQL(module.GJdbcUser, module.GJdbcPasswd, module.GFeEntryHost, module.GFeEntryQueryPort, module.GJdbcDb, queryCMD)
    if err != nil{
        infoMess = fmt.Sprintf("Error in run sql when check BE status: [BeHost = %s, error = %v]", tmpBeHost, err)
        utl.Log("DEBUG", infoMess)
        return beStat, err
    }


    for rows.Next(){
        err = rows.Scan(  &tmpBeStat.BackendId,
                          &tmpBeStat.Cluster,
                          &tmpBeStat.IP,
                          &tmpBeStat.HeartbeatServicePort,
                          &tmpBeStat.BePort,
                          &tmpBeStat.HttpPort,
                          &tmpBeStat.BrpcPort,
                          &tmpBeStat.LastStartTime,
                          &tmpBeStat.LastHeartbeat,
                          &tmpBeStat.Alive,
                          &tmpBeStat.SystemDecommissioned,
                          &tmpBeStat.ClusterDecommissioned,
                          &tmpBeStat.TabletNum,
                          &tmpBeStat.DataUsedCapacity,
                          &tmpBeStat.AvailCapacity,
                          &tmpBeStat.TotalCapacity,
                          &tmpBeStat.UsedPct,
                          &tmpBeStat.MaxDiskUsedPct,
                          &tmpBeStat.ErrMsg,
                          &tmpBeStat.Version,
                          &tmpBeStat.Status,
                          &tmpBeStat.DataTotalCapacity,
                          &tmpBeStat.DataUsedPct)
        if err != nil {
            infoMess = fmt.Sprintf("Error in scan sql result [BeHost = %s, error = %v]", tmpBeHost, err)
            utl.Log("DEBUG", infoMess)
            return beStat, err
        }

        if string(tmpBeStat.IP) == tmpBeHost && tmpBeStat.HeartbeatServicePort == tmpHeartbeatServicePort {
            beStat = tmpBeStat
            //GFeStatusArr[feId] = feStat
            return beStat, nil
        }
    }

    return beStat, err
}


func CheckBeStatus(beId int) (beStat BeStatusStruct, err error) {

    var bePortRun   bool
    bePortRun, err = CheckBePortStatus(beId)

    if bePortRun {
        beStat, err = GetBeStatJDBC(beId)
    }

    return beStat, err
}



