import { Task, ViewTask } from '../models/tasks'
import { axiosPost, axiosPut, axiosGet, axiosPatch, axiosDelete } from './axios'
import { taskToViewTask, viewTaskToTask } from './convert'

const API_PREFIX = "/api/v1"

// export async function AddTask(task: ViewTask, fn: (ret: ViewTask| false)=>(void)) {
//   let res = await axiosPost(`${API_PREFIX}/tasks`, task, true)
//   if (res) {
//     let ret = taskToViewTask(res)
//     fn(ret)
//   }
//   fn(res)
// }

// export async function QueryAllTask(fn: (data: ViewTask[]|false) => (void)) {
//   let res: any[] | false = await axiosGet(`${API_PREFIX}/tasks`)
//   let ret = res
//   if (res) {
//     ret = res.map((one: Task) => taskToViewTask(one))
//   }
//   fn(ret)
// }

// export async function UpdateTask(task: ViewTask, fn: (success: boolean)=>(void)) {
//   let data = viewTaskToTask(task)
//   let res = await axiosPatch(`${API_PREFIX}/tasks`, data)
//   fn(res)
// }

// export async function DeleteTask(id: number, fn: (success: boolean)=>(void)) {
//   let res = await axiosDelete(`${API_PREFIX}/tasks/${id}`)
//   fn(res)
// }

// export async function ChangePosition(id: number, position: number, fn: (success: boolean)=>(void)) {
//   let res = await axiosPut(`${API_PREFIX}/tasks/${id}/position/${position}`, {})
//   fn(res)
// }

export type MilvusCluster = {
  name: string;
  namespace: string;
  status: string;
  version: string
  managed_by: string
};

export type ClassifiedPods = {
  components: string[]
  component_pods: Map<string, string[]>
};

export type ClassifiedPodsRaw = {
  components: string[]
  component_pods: any
};

export type Response<T> = {
  code: number;
  message: string;
  data: T;
}

export async function ListMilvusCluster(cluster: string, fn: (data: MilvusCluster[]) => (void)) {
  let res: Response<MilvusCluster[]> | false
  if (isDev()) {
    console.info('ListMilvusCluster', cluster)
    fn([
      {
        name: 'test1',
        namespace: "ns1",
      },
      {
        name: 'test2',
        namespace: "ns2",
      },
    ] as MilvusCluster[])
    return
  }
  res = await axiosGet(`${API_PREFIX}/clusters/${cluster}/milvus`)
  if (res) {
    fn(res.data)
  } else {
    fn([])
  }
}

export async function ListNamespaces(cluster: string, fn: (data: string[]) => (void)) {
  let res: Response<string[]> | false
  if (isDev()) {
    console.info('ListNamespaces', cluster)
    fn([
      "ns1","ns2"
    ])
    return
  }
  res = await axiosGet(`${API_PREFIX}/clusters/${cluster}/milvus`)
  if (res) {
    fn(res.data)
  } else {
    fn([])
  }
}

export async function ListClassfiedPods(cluster: string, ns: string, milvus: string, by: string, fn: (data: ClassifiedPods) => (void)) {
  let res: Response<ClassifiedPodsRaw> | false
  if (isDev()) {
    console.info('ListNamespaces', cluster)
    fn({} as ClassifiedPods)
    return
  }
  res = await axiosGet(`${API_PREFIX}/clusters/${cluster}/milvus/${ns}/${milvus}/pods?by=${by}`)
  if (res) {
    let data = res.data
    let map = new Map<string, string[]>()
    data.components.map((component: string) => {
      return map.set(component, data.component_pods[component])
    })
    fn({components: data.components, component_pods: map})
  } else {
    fn({} as ClassifiedPods)
  }
}

function isDev(): boolean {
  return process.env.NODE_ENV === 'development' && false
}