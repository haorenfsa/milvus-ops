import moment from "moment";

// common
export enum TaskStatus {
  TODO,
  Doing,
  Done,
  Pending,
  Closed
}

export function GetEnumKeys(enumType: Object): string[] {
  let ret = [] as string[]
  for (let item in enumType) {
    if (isNaN(Number(item))) {
        ret.push(item)
    }
  }
  return ret
}

// common
export interface Task {
  id: number;
  name: string;
  status: TaskStatus;
  plan: TaskPlan
  project: string;
}

export interface TaskPlan {
  year: number
  month: number
  week: number
  day: number
}

export interface TaskPlanView {
  moment: moment.Moment | null
  level: string
}

// view
export interface ViewTask {
  id: number;
  name: string;
  plan: TaskPlanView
  status: TaskStatus;
  project: string;
  editing?: boolean;
  editingName?: string;
}
