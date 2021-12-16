import { ViewTask, Task, TaskPlanView, TaskPlan } from "../models/tasks"
import moment from "moment"

// convert utility
export function dateToString(date: Date): string {
  return date.toISOString()
}

export function viewTaskToTask(task: ViewTask): Task {
  let ret = {
    id: task.id,
    name: task.name,
    plan: taskPlanFromView(task.plan),
    status: task.status,
    project: task.project
  }
  return ret
}

export function taskToViewTask(task: Task): ViewTask {
  let ret = {
    id: task.id,
    name: task.name,
    plan: taskPlanToView(task.plan),
    status: task.status,
    project: task.project
  }
  console.log(ret.plan)
  return ret
}

export function taskPlanFromView(view: TaskPlanView): TaskPlan {
  let ret = {
    year: -1,
    month: -1,
    week: -1,
    day: -1,
  }
  let moment = view.moment
  let level = view.level
  if (!moment) {
    return ret
  }

  ret.year = moment.year()
  if (level === "year") {
    return ret
  }
  ret.month = moment.month()
  if (level === "month") {
    return ret
  }
  ret.week = moment.week()
  if (level === "week") {
    return ret
  }
  ret.day = moment.day()
  return ret
}

export function taskPlanToView(plan: TaskPlan): TaskPlanView {
  let ret = {
    level: "none"
  } as TaskPlanView
  if (plan.day > 0) {
    ret.moment = moment(`${plan.year}-W${plan.week}-${plan.day}`)
    ret.level = "day"
    return ret
  }

  if (plan.week > 0) {
    ret.moment = moment(`${plan.year}-W${plan.week}`)
    ret.level = "week"
    return ret
  }

  if (plan.month > 0) {
    ret.moment = moment(`${plan.year}-${plan.month}`)
    ret.level = "month"
    return ret
  }

  if (plan.year > 0) {
    ret.moment = moment(`${plan.year}`)
    ret.level = "year"
    return ret
  }
  return ret
}