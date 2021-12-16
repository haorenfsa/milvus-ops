import React from "react";
import { Table } from "antd";
import { DndProvider, DragSource, DropTarget } from "react-dnd";
import HTML5Backend from "react-dnd-html5-backend";
import update from "immutability-helper";
import { TableProps } from "antd/lib/table";
import './DragableTable.css'

let dragingIndex = -1;
class BodyRow extends React.Component<any, any> {
  render() {
    const {
      isOver,
      connectDragSource,
      connectDropTarget,
      moveRow,
      ...restProps
    } = this.props;
    const style = { ...restProps.style, cursor: "move" };
    let { className } = restProps;
    if (isOver) {
      if (restProps.index > dragingIndex) {
        className += " drop-over-downward";
      }
      if (restProps.index < dragingIndex) {
        className += " drop-over-upward";
      }
    }

    return connectDragSource(
      connectDropTarget(
        <tr {...restProps} className={className} style={style} />
      )
    );
  }
}

const rowSource = {
  beginDrag(props: any) {
    dragingIndex = props.index;
    return {
      index: props.index,
    };
  },
};

const rowTarget = {
  drop(props: any, monitor: any) {
    const dragIndex = monitor.getItem().index;
    const hoverIndex = props.index;

    // Don't replace items with themselves
    if (dragIndex === hoverIndex) {
      return;
    }

    // Time to actually perform the action
    props.moveRow(dragIndex, hoverIndex);

    // Note: we're mutating the monitor item here!
    // Generally it's better to avoid mutations,
    // but it's good here for the sake of performance
    // to avoid expensive index searches.
    monitor.getItem().index = hoverIndex;
  },
};

const DragableBodyRow = DropTarget("row", rowTarget, (connect, monitor) => ({
  connectDropTarget: connect.dropTarget(),
  isOver: monitor.isOver(),
}))(
  DragSource("row", rowSource, (connect) => ({
    connectDragSource: connect.dragSource(),
  }))(BodyRow)
);

export interface DragTableProps extends TableProps<any> {
  data: any[];
  onMoveRow: (oldIndex: number, newIndex: number) => Promise<boolean>;
}

interface DragTableState {
  data: any[];
}

class DragTable extends React.Component<DragTableProps, DragTableState> {
  state = {
    data: this.props.data,
  };

  components = {
    body: {
      row: DragableBodyRow,
    },
  };

  componentWillReceiveProps(props: any) {
    this.setState({ data: props.data });
  }

  moveRow = (dragIndex: any, hoverIndex: any) => {
    const { data } = this.state;
    const dragRow = data[dragIndex];
    this.props.onMoveRow(dragIndex, hoverIndex).then(() =>
      this.setState(
        update(this.state, {
          data: {
            $splice: [
              [dragIndex, 1],
              [hoverIndex, 0, dragRow],
            ],
          },
        })
      )
    );
  };

  render() {
    return (
      <DndProvider backend={HTML5Backend}>
        <Table
          pagination={false}
          dataSource={this.state.data}
          components={this.components}
          onRow={(record, index) => ({
            index,
            moveRow: this.moveRow,
          })}
          {...this.props}
        />
      </DndProvider>
    );
  }
}

export default DragTable;
