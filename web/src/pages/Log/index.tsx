import { Button, DatePicker, Input, InputNumber, message, Radio, Select, Switch, Tag, TimePicker } from "antd";
import React, { useState, useEffect } from "react";
import { XTerm } from 'xterm-for-react'
import { FitAddon } from 'xterm-addon-fit';
import ButtonGroup from "antd/lib/button/button-group";
import { mustSplit3, parseQuery, queryGet, setQuery } from "../util";
import { ClassifiedPods, ListClassfiedPods, ListMilvusCluster, MilvusCluster } from "../../api/milvus";
import moment from 'moment';

const format = 'HH:mm';


const xtermRef = React.createRef<XTerm>()
let ws: WebSocket

const Log: React.FC = () => {
  let query = parseQuery()
  const [cluster, setCluster] = React.useState("qa")
  const [comPods, setComPods] = React.useState({} as ClassifiedPods)
  const [pods, setPods] = React.useState([] as string[])
  const [milvuses, setMilvuses] = React.useState([] as MilvusCluster[])
  const [milvus, setMilvus] = React.useState(queryGet(query, "ns") + "/" + queryGet(query, "milvus") + "/" + queryGet(query, "by"))
  const [component, setSelectedComponent] = React.useState(queryGet(query, "component"))
  const [pod, setSelectedPod] = React.useState(queryGet(query, "pod"))
  const [container, setSelectedContainer] = React.useState(queryGet(query, "container"))
  const [sinceTime, setSinceTime] = React.useState(moment().toISOString())
  const [limitSizeMB, setLimitSizeMB] = React.useState(1)

  useEffect(() => {
    if (milvus[0] === "" && milvuses.length > 0) {
      setMilvus(milvuses[0].namespace + "/" + milvuses[0].name + "/" + milvuses[0].managed_by)
    } else {
      ListMilvusCluster(cluster, setMilvuses)
    }
  }, [])

  useEffect(() => {
    setQuery(query, "pod", pod)
  }, [pod])

  useEffect(() => {
    setQuery(query, "component", component)
    let pods = comPods.component_pods && comPods.component_pods.get(component)
    if (pods) {
      setPods(pods)
      if (pods.length > 0) {
        setSelectedPod(pods[0])
      }
    }
  }, [component])

  useEffect(() => {
    let splited = mustSplit3(milvus)
    setQuery(query, "ns", splited[0])
    setQuery(query, "milvus", splited[1])
    setQuery(query, "by", splited[2])
    if (milvus !== "") {
      initClassifedPods()
    }
  }, [milvus])

  function initClassifedPods() {
    let splited = mustSplit3(milvus)
    ListClassfiedPods(cluster, splited[0], splited[1], splited[2], (res: ClassifiedPods) => {
      console.log("pods", res)
      setComPods(res)
      if (res.component_pods) {
        let proxys = res.component_pods.get("proxy")
        setSelectedComponent("proxy")
        if (proxys && proxys.length > 0) {
          setSelectedPod(proxys[0])
        }

        let alone = res.component_pods.get("standalone")
        if (alone && alone.length > 0) {
          setSelectedComponent("standalone")
          setSelectedPod(alone[0])
        }
      }
    })
  }

  function connect(xtermRef: React.RefObject<XTerm>) {
    const fitAddon = new FitAddon();
    let term = xtermRef.current!.terminal
    if (term) {
      term.reset()
      term.paste = (data: string) => {
        console.log("paste2")
        term.write(data)
      }

      if (term.element) {
        term.element.addEventListener('paste', (e: any) => {
          console.log("paste", e)
          term.write(e.target.value)
        })
        let ele = document.getElementById('terminal-container')
        let cols = Math.floor((ele!.parentElement!.getBoundingClientRect().width - 200) / 9 - 1)
        if (cols < 40) { cols = 40 }
        term.resize(cols, 40);
      }
      fitAddon.activate(term);
      // fitAddon.fit();
      term.onResize((size: { cols: number, rows: number }) => {
        const cols = size.cols;
        const rows = size.rows;
        let msg = JSON.stringify({
          operation: "resize",
          row: rows,
          col: cols
        })
        console.log("resize", msg)
        ws.send(msg)
      });
    }

    let host = window.location.host
    if (host.indexOf("localhost") !== -1) {
      host = "localhost:8080"
    }
    console.log("connecting to websocket ", host)
    ws && ws.close()
    let splited = mustSplit3(milvus)
    ws = new WebSocket(`ws://${host}/api/v1/clusters/${cluster}/milvus/${splited[0]}/${splited[1]}/logs?pod=${pod}&container=${container}&component=${component}?by=${splited[2]}`)
    ws.onopen = () => {
      term.write("connecting...\n")
      ws.send(JSON.stringify({
        operation: "resize",
        row: term?.rows,
        col: term?.cols
      }))
    }
    ws.onerror = (e) => {
      message.warning("connection failed")
    }
    ws.onclose = (e) => {
      message.warning("connection closed")
    }
    ws.onmessage = (e) => {
      if (xtermRef.current) {
        let data = JSON.parse(e.data)
        if (data.operation === "stdout") {
          xtermRef.current.terminal.write(data.data)
        }
      }
    }
    return () => {
      ws.close()
    }
  }

  function onKey(e: { key: string, domEvent: KeyboardEvent }) {
    const ev = e.domEvent;
    if (e.key.charCodeAt(0) === 13) {
      xtermRef?.current?.terminal.writeln(e.key);
    } else {
      xtermRef?.current?.terminal.write(e.key);
    }
  }

  let splited = mustSplit3(milvus)
  let podSelection = pods.map((podName) => <Select.Option key={`pod-${podName}`} value={podName}>
      {podName.substring(splited[1].length)}
    </Select.Option>)

  let scheme = window.location.protocol === "https" ? window.location.protocol : "http"
  let host = window.location.host
  if (host.indexOf("localhost") !== -1) {
    host = "localhost:8080"
  }

  return (
    <div>
      <h1>Logs</h1>
      <div style={{ marginBottom: 8 }}>
        Milvus: <span style={{ marginRight: 8 }} />
        <Select onChange={setMilvus} showSearch style={{ width: 400 }} value={milvus}>
          {milvuses.map((milvus) => <Select.Option key={`milvus-${milvus.name}`} value={milvus.namespace + "/" + milvus.name + "/" + milvus.managed_by}><Tag style={{ textAlign: "right" }}>{milvus.namespace + "/"}</Tag><Tag color="green">{milvus.name}</Tag></Select.Option>)}
        </Select>
        <span style={{ marginRight: 16 }} />
        Component: <span style={{ marginRight: 8 }} />
        <Select showSearch style={{ width: 150 }} onChange={setSelectedComponent} value={component}>
          {comPods.components && comPods.components.map((c) => {
            return <Select.Option key={`com-${c}`} value={c}>{c}</Select.Option>
          })}
        </Select>
        <span style={{ marginRight: 16 }} />
        Pod: <span style={{ marginRight: 8 }} />
        <Select showSearch value={pod} onChange={setSelectedPod} style={{ width: 300 }}>
          {podSelection}
        </Select>
        <span style={{ marginRight: 16 }} />
      </div>
      <div style={{ marginBottom: 8 }}>
        Download History Log: Since <span style={{ marginRight: 8 }} /> 
          <DatePicker value={moment(sinceTime)} onChange={(dt) => dt && setSinceTime(dt.toISOString())} showTime placeholder="select date time"/> <span style={{ marginRight: 8 }} />
          SizeLimit <span style={{ marginRight: 8 }} />
          <InputNumber onChange={(val) => val && setLimitSizeMB(val)} defaultValue={limitSizeMB}/> MB <span style={{ marginRight: 8 }} />
          <Button href={`${scheme}://${host}/api/v1/clusters/${cluster}/milvus/${splited[0]}/${splited[1]}/files/log?pod=${pod}&container=${container}&component=${component}?by=${splited[2]}&since=${sinceTime}&size_mb=${limitSizeMB}`} target="_blank" type="primary" shape="round" icon="download">Download</Button>
      </div>
      <div style={{ marginBottom: 8 }}>
          Follow log: <span style={{ marginRight: 8 }} />
          <ButtonGroup>
            <Button onClick={() => connect(xtermRef)} type="primary">Start </Button>
            <Button onClick={() => ws.close()} type="danger">Stop </Button>
          </ButtonGroup>
        </div>
      <div id="terminal-container">
        <XTerm options={{
          rendererType: 'canvas',
          cursorBlink: true,
          convertEol: true,
          cursorStyle: 'bar',
          scrollback: 10000,
          rows: 40,
          cols: 128,
          theme: theme,
        }} key="the-term" ref={xtermRef} onKey={onKey} />
      </div>
    </div>
  );
}

const theme = {
  foreground: "#ffffff",
  background: "#1b212f",
  cursor: "#ffffff",
  selection: "rgba(255, 255, 255, 0.3)",
  black: "#000000",
  brightBlack: "#808080",
  red: "#ce2f2b",
  brightRed: "#f44a47",
  green: "#00b976",
  brightGreen: "#05d289",
  yellow: "#e0d500",
  brightYellow: "#f4f628",
  magenta: "#bd37bc",
  brightMagenta: "#d86cd8",
  blue: "#1d6fca",
  brightBlue: "#358bed",
  cyan: "#00a8cf",
  brightCyan: "#19b8dd",
  white: "#e5e5e5",
  brightWhite: "#ffffff"
}

export default Log;