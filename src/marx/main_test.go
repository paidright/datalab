package main

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarx(t *testing.T) {
	headers := strings.Split("employee_id,PayPeriod_End_Date,payperiod_WoWID,payperiod_id,Last_day_at_work,First_Day_At_Work,dob,Position,Position_Name,Company,Group,Brand,OPS_Support,Region,Area,Location,Paying_Department,Personnel_Area_Sub_Area,Type_of_Employee_Movement_Label,Employee_Movement_End_Date,Employee_Group_Label,Employee_Subgroup_Label,Kronos_Employee_Label,Employee_Status_Label,Working_Days_Per_Week,Base_Hours,Leave_Entitlement_Label,Short_Term_Incentive_Plan,Car_Eligibility_Label,Pay_Scale_Type_Label,Pay_Scale_Area_Label,Pay_Scale_Group_code,Pay_Scale_Level_code,Superannuation_Guarantee,Band_Band,gender,Line_Manager", ",")
	result := strings.Builder{}
	output := csv.NewWriter(&result)
	err := processFile("./test_data/input.csv", headers, output)
	assert.Nil(t, err)
}

func BenchmarkMarx(b *testing.B) {
	headers := strings.Split("employee_id,PayPeriod_End_Date,payperiod_WoWID,payperiod_id,Last_day_at_work,First_Day_At_Work,dob,Position,Position_Name,Company,Group,Brand,OPS_Support,Region,Area,Location,Paying_Department,Personnel_Area_Sub_Area,Type_of_Employee_Movement_Label,Employee_Movement_End_Date,Employee_Group_Label,Employee_Subgroup_Label,Kronos_Employee_Label,Employee_Status_Label,Working_Days_Per_Week,Base_Hours,Leave_Entitlement_Label,Short_Term_Incentive_Plan,Car_Eligibility_Label,Pay_Scale_Type_Label,Pay_Scale_Area_Label,Pay_Scale_Group_code,Pay_Scale_Level_code,Superannuation_Guarantee,Band_Band,gender,Line_Manager", ",")
	for n := 0; n < b.N; n++ {
		result := strings.Builder{}
		output := csv.NewWriter(&result)
		processFile("./test_data/input.csv", headers, output)
	}
}
