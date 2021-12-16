import { Button, message, Radio, Select, Tag } from "antd";
import React, { useState, useEffect, useRef, useCallback } from "react";
import { XTerm } from 'xterm-for-react'
import { FitAddon } from 'xterm-addon-fit';
import { addListener } from "process";
import { mustSplit3, parseQuery, queryGet, setQuery } from "../util";
import { ClassifiedPods, ListClassfiedPods, ListMilvusCluster, MilvusCluster } from "../../api/tasks";
import { Terminal } from 'xterm'

const xtermRef = React.createRef<XTerm>()
let ws: WebSocket

const Shell: React.FC = () => {
  let query = parseQuery()
  const [cluster, setCluster] = React.useState("qa")
  const [comPods, setComPods] = React.useState({} as ClassifiedPods)
  const [pods, setPods] = React.useState([] as string[])
  const [milvuses, setMilvuses] = React.useState([] as MilvusCluster[])
  const [milvus, setMilvus] = React.useState(queryGet(query, "ns") + "/" + queryGet(query, "milvus") + "/" + queryGet(query, "by"))
  const [component, setSelectedComponent] = React.useState(queryGet(query, "component"))
  const [pod, setSelectedPod] = React.useState(queryGet(query, "pod"))
  const [container, setSelectedContainer] = React.useState(queryGet(query, "container"))

  useEffect(() => {
    if (milvus[0] === "" && milvuses.length > 0) {
      setMilvus(milvuses[0].namespace + "/" + milvuses[0].name + "/" + milvuses[0].managed_by)
    } else {
      ListMilvusCluster(cluster, setMilvuses)
    }
  }, [])

  useEffect(() => {
    setQuery(query, "pod", pod)
    if (pod !== "") {
      connect(xtermRef)
    }
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
    ws = new WebSocket(`ws://${host}/api/v1/clusters/${cluster}/milvus/${splited[0]}/${splited[1]}/shell?pod=${pod}&container=${container}&component=${component}?by=${splited[2]}`)
    ws.onopen = () => {
      term.write("login...\n")
      ws.send(JSON.stringify({
        operation: "resize",
        row: term?.rows,
        col: term?.cols
      }))
    }
    ws.onerror = (e) => {
      message.warning("connection failed")
    }
    ws.onmessage = (e) => {
      if (xtermRef.current) {
        console.log("onmessage", e.data)
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
    // const printable = !ev.altKey && !ev.ctrlKey && !ev.metaKey;

    ws.send(JSON.stringify({
      operation: "stdin",
      data: e.key
    }))
  }

  let splited = mustSplit3(milvus)
  let podSelection = pods.map((podName) => <Select.Option key={`pod-${podName}`} value={podName}>
      {podName.substring(splited[1].length)}
    </Select.Option>)

  return (
    <div>
      <h1>Web Shell</h1>
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
        {/* <span style={{ marginRight: 8 }}>Namespace:</span>
        <Select showSearch disabled style={{ width: 100 }} defaultValue={"all"} />
        <span style={{ marginRight: 16 }} /> */}
        K8s Cluster: <span style={{ marginRight: 8 }} />
        <Radio.Group disabled onChange={(e) => setCluster(e.target.value)} defaultValue={cluster}>
          <Radio.Button value="qa">QA</Radio.Button>
          <Radio.Button value="ci">CI</Radio.Button>
        </Radio.Group>
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

export default Shell;