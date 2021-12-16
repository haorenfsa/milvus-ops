import React from 'react'
import { Select, DatePicker } from 'antd'
import { DatePickerProps } from 'antd/lib/date-picker/interface'

interface PlanPickerState {
  pickerLevel: string
}

interface PlanPickerProps extends DatePickerProps {
  onChangePicker?: (moment: any, pickerLevel: string) => void
  defaultPickerLevel: string
}

class PlanPicker extends React.Component<PlanPickerProps, PlanPickerState> {
  state = {
    pickerLevel: this.props.defaultPickerLevel,
  }

  onChange = () => (moment: any) => {
    if (this.props.onChangePicker) {
      this.props.onChangePicker(moment, this.state.pickerLevel)
    }
  }

  render() {
    let { pickerLevel } = this.state
    if (pickerLevel === "none") {
      pickerLevel = "week"
    }
    return (
      <div>
        <Select style={{width: 75}} onChange={(value: string)=> this.setState({pickerLevel: value})} defaultValue={pickerLevel}>
          <Select.Option key="day">day</Select.Option>
          <Select.Option key="week">week</Select.Option>
          <Select.Option key="month">month</Select.Option>
        </Select>
        {
          pickerLevel === "day" && <DatePicker style={{width: 200}} onChange={this.onChange()} {...this.props}/>
        }
        {
          pickerLevel === "week" && <DatePicker.WeekPicker onChange={this.onChange()} {...this.props}/>
        }
        {
          pickerLevel === "month" && <DatePicker.MonthPicker onChange={this.onChange()} {...this.props}/>
        }
        
      </div>
    )
  }
}
export default PlanPicker
