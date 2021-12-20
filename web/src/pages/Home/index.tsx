import React, { Fragment } from "react";
import { Input, Button, Row, Col, Select, Icon, Radio, Table } from "antd";

import {
  ListMilvusCluster,
  ListNamespaces,
  MilvusCluster,
} from "../../api/tasks";
import { Pagination } from "../../models/page";

import PlanPicker from "../../components/PlanPicker";
import DragTable from "../../components/DragableTable";
import ButtonGroup from "antd/lib/button/button-group";


const selectAll = "_all"

const Welcome: React.FC = () => {
  // const intl = useIntl();
  const [cluster, setCluster] = React.useState("qa")
  const [mcList, setMcList] = React.useState([] as MilvusCluster[])
  const [nss, setNss] = React.useState([] as string[]) // namespaces
  const [loading, setLoading] = React.useState(true)
  const [search = "", setSearch] = React.useState("" as string)

  React.useEffect(() => {
    ListNamespaces(cluster, setNss)
    ListMilvusCluster(cluster, (mcList) => {setMcList(mcList); setLoading(false)})
  }, [cluster])

  let filterMcList = () => {
    if (search === "") return mcList
    return mcList.filter(mc => mc.name.includes(search))
  }
  let filteredList = filterMcList()

  const columns = [
    {
      title: 'Milvus Name',
      dataIndex: 'name',
      width: 300,
    },
    {
      title: 'Operation',
      dataIndex: 'operation',
      render: (text: any, record: MilvusCluster) => {
        return (
          <div>
            <ButtonGroup>
              <Button href={`/app/logs?cluster=${cluster}&ns=${record.namespace}&milvus=${record.name}&by=${record.managed_by}`} type="primary">log</Button>
              <Button href={`/app/shell?cluster=${cluster}&ns=${record.namespace}&milvus=${record.name}&by=${record.managed_by}`} type="primary">login</Button>
            </ButtonGroup>
          </div>
        )
      }
    },
    {
      title: 'Namespace',
      dataIndex: 'namespace',
    },
    {
      title: 'Status',
      dataIndex: 'status',
    },
    {
      title: 'ManagedBy',
      dataIndex: 'managed_by',
    },
    {
      title: 'Version',
      dataIndex: 'version',
    },
  ]

  return (
    <div>
      <div>
        <h3>Milvus</h3>
        <Input.Search autoFocus style={{width: 300}} placeholder="Search by name" onChange={(e: any) => {setSearch(e.target.value)}} enterButton />
        <span style={{ marginRight: 16 }} />
        <span style={{ marginRight: 8 }}>Namespace:</span>
        <Select showSearch disabled style={{ width: 100 }} defaultValue={selectAll}>
          <Select.Option key={"ns-all"} value={selectAll}>ALL</Select.Option>
          {/* {nss.map((ns) => <Select.Option key={"ns-"+ns} value={ns}>{ns}</Select.Option>)} */}
        </Select>
        <span style={{ marginRight: 16 }} />

        K8s Cluster: <span style={{ marginRight: 8 }} />
        <Radio.Group onChange={(e) => setCluster(e.target.value)} defaultValue={cluster}>
          <Radio.Button value="qa">QA</Radio.Button>
          <Radio.Button value="ci">CI</Radio.Button>
        </Radio.Group>
      </div>
      <div style={{ marginTop: 16 }} />
      <Table pagination={{pageSize: 20}} loading={loading} rowKey={(record) => `mc-${record.namespace}-${record.name}`} columns={columns} dataSource={filteredList} />
    </div>
  );
};

export default Welcome;