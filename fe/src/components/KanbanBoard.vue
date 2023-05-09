<template>
    <div ref="myDiagramDivRef" style="width: 100%; height: 100%; position: relative;"></div>
  </template>
  
  <script lang="ts" setup>
  import * as go from 'gojs';
  import { ref, onMounted } from 'vue';
  
  const $ = go.GraphObject.make;
  const myDiagramDivRef = ref()
  var myDiagram: go.Diagram
  var myModel = $(go.GraphLinksModel);
  
  
  class PoolLayout extends go.GridLayout {
    private MINLENGTH: number = 200;
    private MINBREADTH: number = 100;
  
    constructor() {
      super();
      this.cellSize = new go.Size(1, 1);
      this.wrappingColumn = Infinity;
      this.wrappingWidth = Infinity;
      this.spacing = new go.Size(0, 0);
      this.alignment = go.GridLayout.Position;
    }
  
    doLayout(coll: any) {
      const diagram = this.diagram;
      if (diagram === null) return;
      diagram.startTransaction("PoolLayout");
      const minlen = this.computeMinPoolLength();
      diagram.findTopLevelGroups().each(lane => {
        if (!(lane instanceof go.Group)) return;
        const shape = lane.selectionObject;
        if (shape !== null) {
          const sz = this.computeLaneSize(lane);
          shape.width = (!isNaN(shape.width)) ? Math.max(shape.width, sz.width) : sz.width;
          shape.height = minlen;
          const cell = lane.resizeCellSize;
          if (!isNaN(shape.width) && !isNaN(cell.width) && cell.width > 0) shape.width = Math.ceil(shape.width / cell.width) * cell.width;
          if (!isNaN(shape.height) && !isNaN(cell.height) && cell.height > 0) shape.height = Math.ceil(shape.height / cell.height) * cell.height;
        }
      });
      super.doLayout(coll);
      diagram.commitTransaction("PoolLayout");
    };
  
    computeMinPoolLength() {
      let len = this.MINLENGTH;
      myDiagram.findTopLevelGroups().each(lane => {
        const holder = lane.placeholder;
        if (holder !== null) {
          const sz = holder.actualBounds;
          len = Math.max(len, sz.height);
        }
      });
      return len;
    }
  
    computeLaneSize(lane: go.Group) {
      const sz = new go.Size(lane.isSubGraphExpanded ? this.MINBREADTH : 1, this.MINLENGTH);
      if (lane.isSubGraphExpanded) {
        const holder = lane.placeholder;
        if (holder !== null) {
          const hsz = holder.actualBounds;
          sz.width = Math.max(sz.width, hsz.width);
        }
      }
      const hdr = lane.findObject("HEADER");
      if (hdr !== null) sz.width = Math.max(sz.width, hdr.actualBounds.width);
      return sz;
    }
  }
  
  function init() {
    myDiagram =
      $(go.Diagram, myDiagramDivRef.value,
        {
          contentAlignment: go.Spot.TopLeft,
          layout: $(PoolLayout),
          mouseDrop: (e: go.InputEvent) => {
            e.diagram.currentTool.doCancel();
          },
          "commandHandler.copiesGroupKey": true,
          "SelectionMoved": relayoutDiagram,
          "SelectionCopied": relayoutDiagram,
          "undoManager.isEnabled": true,
          "textEditingTool.starting": go.TextEditingTool.SingleClick
        });
  
    // myDiagram.toolManager.draggingTool.doActivate = function () {  // method override must be function, not =>
    //   go.DraggingTool.prototype.doActivate.call(this);
    //   if (this.currentPart) {
    //     this.currentPart.opacity = 0.6;
    //     this.currentPart.layerName = "Foreground";
    //   }
    // }
    // myDiagram.toolManager.draggingTool.doDeactivate = function () {  // method override must be function, not =>
    //   if (this.currentPart) {
    //     this.currentPart.opacity = 1;
    //     this.currentPart.layerName = "";
    //   }
    //   go.DraggingTool.prototype.doDeactivate.call(this);
    // }
  
    function relayoutDiagram() {
      myDiagram.selection.each(n => n.invalidateLayout());
      myDiagram.layoutDiagram();
    }
  
    const noteColors = ['#009CCC', '#CC293D', '#FFD700'];
    function getNoteColor(num: number) {
      return noteColors[Math.min(num, noteColors.length - 1)];
    }
  
    myDiagram.nodeTemplate =
      $(go.Node, "Horizontal",
        new go.Binding("location", "loc", go.Point.parse).makeTwoWay(go.Point.stringify),
        $(go.Shape, "Rectangle", {
          fill: '#009CCC', strokeWidth: 1, stroke: '#009CCC',
          width: 6, stretch: go.GraphObject.Vertical, alignment: go.Spot.Left,
          click: (e: go.InputEvent, obj: go.GraphObject) => {
            myDiagram.startTransaction("Update node color");
            let newColor = parseInt(obj.part?.data.color) + 1;
            if (newColor > noteColors.length - 1) newColor = 0;
            myDiagram.model.setDataProperty(obj.part?.data, "color", newColor);
            myDiagram.commitTransaction("Update node color");
          }
        },
          new go.Binding("fill", "color", getNoteColor),
          new go.Binding("stroke", "color", getNoteColor)
        ),
        $(go.Panel, "Auto",
          $(go.Shape, "Rectangle", { fill: "white", stroke: '#CCCCCC' }),
          $(go.Panel, "Table",
            { width: 130, minSize: new go.Size(NaN, 50) },
            $(go.TextBlock,
              {
                name: 'TEXT',
                margin: 6, font: '11px Lato, sans-serif', editable: true,
                stroke: "#000", maxSize: new go.Size(130, NaN),
                alignment: go.Spot.TopLeft
              },
              new go.Binding("text", "text").makeTwoWay())
          )
        )
      );
  
    function highlightGroup(grp: go.Group, show: boolean) {
      const draggingTool = myDiagram.toolManager.draggingTool;
      if (show) {
        const part = draggingTool.currentPart;
        if (part instanceof go.Part && part.containingGroup !== grp) {
          grp.isHighlighted = true;
          return;
        }
      }
      grp.isHighlighted = false;
    }
  
    myDiagram.groupTemplate =
      $(go.Group, "Vertical",
        {
          selectable: true,
          selectionObjectName: "SHAPE",
          layerName: "Background",
          movable: true,
          layout: $(go.GridLayout,
            {
              wrappingColumn: 1,
              cellSize: new go.Size(1, 1),
              spacing: new go.Size(5, 5),
              alignment: go.GridLayout.Position,
              comparer: (a: go.Part, b: go.Part) => {  // can re-order tasks within a lane
                const ay = a.location.y;
                const by = b.location.y;
                if (isNaN(ay) || isNaN(by)) return 0;
                if (ay < by) return -1;
                if (ay > by) return 1;
                return 0;
              }
            }),
          click: (e, grp) => {
            if (!e.shift && !e.control && !e.meta) e.diagram.clearSelection();
          },
          computesBoundsAfterDrag: true,
          handlesDragDropForMembers: true,
          mouseDragEnter: (e, grp: go.GraphObject, prev) => highlightGroup(grp as go.Group, true),
          mouseDragLeave: (e, grp: go.GraphObject, next) => highlightGroup(grp as go.Group, false),
          mouseDrop: (e, grp: go.GraphObject) => {
            if (e.diagram.selection.all(n => !(n instanceof go.Group))) {
              if (grp.diagram) {
                const ok = (grp as go.Group).addMembers(grp.diagram.selection, true);
                if (!ok) grp.diagram.currentTool.doCancel();
              }
            }
          },
          subGraphExpandedChanged: (grp: go.Group) => {
            const shp = grp.selectionObject;
            if (grp.diagram?.undoManager.isUndoingRedoing) return;
            if (grp.isSubGraphExpanded) {
              shp.width = grp.data.savedBreadth;
            } else {
              if (!isNaN(shp.width)) grp.diagram?.model.set(grp.data, "savedBreadth", shp.width);
              shp.width = NaN;
            }
          }
        },
        new go.Binding("location", "loc", go.Point.parse).makeTwoWay(go.Point.stringify),
        new go.Binding("isSubGraphExpanded", "expanded").makeTwoWay(),
        $(go.Panel, "Horizontal",
          { name: "HEADER", alignment: go.Spot.Left },
          // $("SubGraphExpanderButton", { margin: 5 }),
          $(go.TextBlock,
            { font: "15px Lato, sans-serif", editable: true, margin: new go.Margin(2, 0, 0, 14) },
            new go.Binding("visible", "isSubGraphExpanded").ofObject(),
            new go.Binding("text").makeTwoWay())
        ),
        $(go.Panel, "Auto",
          $(go.Shape, "Rectangle",
            {
              name: "SHAPE", fill: "#F1F1F1", stroke: null, strokeWidth: 4, width: 162, click: (e: any, node: go.GraphObject) => {
                console.log(e)
                console.log(node.width)
              }
            },  // strokeWidth controls space between lanes
            new go.Binding("fill", "isHighlighted", h => h ? "#D6D6D6" : "#B6C4D3").ofObject(),
            new go.Binding("desiredSize", "size", go.Size.parse).makeTwoWay(go.Size.stringify)),
          $(go.Placeholder,
            { padding: 12, alignment: go.Spot.TopLeft }),
          $(go.TextBlock,
            {
              name: "LABEL", font: "15px Lato, sans-serif", editable: true,
              angle: 90, alignment: go.Spot.TopLeft, margin: new go.Margin(4, 0, 0, 2)
            },
            new go.Binding("visible", "isSubGraphExpanded", e => !e).ofObject(),
            new go.Binding("text").makeTwoWay())
        ),
        $(go.Panel, "Horizontal",
          {
            row: 2,
            click: (e, node) => {
              e.diagram.startTransaction('add node');
              let sel = e.diagram.selection.first();
              if (!sel) sel = e.diagram.findTopLevelGroups().first();
              if (!sel) return;
              if (!(sel instanceof go.Group)) sel = sel.containingGroup;
              if (!sel) return;
              const newdata = { group: sel.key, loc: "0 9999", text: "New item " + (sel as go.Group).memberParts.count, color: 0 };
              e.diagram.model.addNodeData(newdata);
              e.diagram.commitTransaction('add node');
              const newnode = myDiagram.findNodeForData(newdata);
              e.diagram.select(newnode);
              e.diagram.commandHandler.editTextBlock();
              e.diagram.commandHandler.scrollToPart(newnode as go.Part);
            },
            background: 'white',
            margin: new go.Margin(10, 4, 4, 4)
          },
          $(go.Panel, "Auto",
            $(go.Shape, "Rectangle", { strokeWidth: 0, stroke: null, fill: '#1D3F75' }),
            $(go.Shape, "PlusLine", { margin: 6, strokeWidth: 2, width: 12, height: 12, stroke: 'white', background: '#1D3F75' })
          ),
          $(go.TextBlock, "添加步骤", { font: '10px Lato, sans-serif', margin: 6, })
        )
      );
  
    myDiagram.add(
      $(go.Part, "Table",
        { position: new go.Point(10, 10), selectable: false },
        // $(go.TextBlock, "Key",
        //   { row: 0, font: "700 14px Droid Serif, sans-serif" }),  // end row 0
        // $(go.Panel, "Horizontal",
        //   { row: 1, alignment: go.Spot.Left },
        //   $(go.Shape, "Rectangle",
        //     { desiredSize: new go.Size(10, 10), fill: '#CC293D', margin: 5 }),
        //   $(go.TextBlock, "Halted",
        //     { font: "700 13px Droid Serif, sans-serif" })
        // ),  // end row 1
        // $(go.Panel, "Horizontal",
        //   { row: 2, alignment: go.Spot.Left },
        //   $(go.Shape, "Rectangle",
        //     { desiredSize: new go.Size(10, 10), fill: '#FFD700', margin: 5 }),
        //   $(go.TextBlock, "In Progress",
        //     { font: "700 13px Droid Serif, sans-serif" })
        // ),  // end row 2
        // $(go.Panel, "Horizontal",
        //   { row: 3, alignment: go.Spot.Left },
        //   $(go.Shape, "Rectangle",
        //     { desiredSize: new go.Size(10, 10), fill: '#009CCC', margin: 5 }),
        //   $(go.TextBlock, "Completed",
        //     { font: "700 13px Droid Serif, sans-serif" })
        // ),  // end row 3
        $(go.Panel, "Horizontal",
          {
            row: 1,
            click: (e, node) => {
              // e.diagram.startTransaction('add group');
              myModel.addNodeData({ key: "NewGroup", text: "New Group", isGroup: true, loc: "0 9999" });
  
            },
            background: 'white',
            margin: new go.Margin(10, 4, 4, 4)
          },
          $(go.Panel, "Auto",
            $(go.Shape, "Rectangle", { strokeWidth: 0, stroke: null, fill: '#1D3F75' }),
            $(go.Shape, "PlusLine", { margin: 6, strokeWidth: 2, width: 12, height: 12, stroke: 'white', background: '#1D3F75' })
          ),
          $(go.TextBlock, "添加阶段", { font: '10px Lato, sans-serif', margin: 6, })
        ),
      )
    );
    load();
  }
  
  function load() {
    myModel.nodeDataArray = [
      { key: "GitPull", text: "拉取代码", isGroup: true, loc: "0 23.52284749830794" },
      { key: "Build", text: "编译", isGroup: true, color: "0", loc: "109 23.52284749830794" },
      { key: "PushImage", text: "上传镜像", isGroup: true, color: "0", loc: "235 23.52284749830794" },
      { key: "Publish", text: "发布", isGroup: true, color: "0", loc: "343 23.52284749830794" },
      { key: "Reviewing", text: "Reviewing", isGroup: true, color: "0", loc: "451 23.52284749830794" },
      { key: "Testing", text: "Testing", isGroup: true, color: "0", loc: "562 23.52284749830794" },
      { key: "Customer", text: "Customer", isGroup: true, color: "0", loc: "671 23.52284749830794" },
      { key: "Customer1", text: "Customer1", isGroup: true, color: "0", loc: "780 23.52284749830794" },
      { key: 1, text: "text for oneA", group: "GitPull", color: "0", loc: "12 35.52284749830794" },
      { key: 2, text: "text for oneB", group: "GitPull", color: "1", loc: "12 65.52284749830794" },
      { key: 3, text: "text for oneC", group: "GitPull", color: "0", loc: "12 95.52284749830794" },
      { key: 4, text: "text for oneD", group: "GitPull", color: "1", loc: "12 125.52284749830794" },
      { key: 5, text: "text for twoA", group: "Build", color: "1", loc: "121 35.52284749830794" },
      { key: 6, text: "text for twoB", group: "Build", color: "1", loc: "121 65.52284749830794" },
      { key: 7, text: "text for twoC", group: "PushImage", color: "0", loc: "247 35.52284749830794" },
      { key: 8, text: "text for twoD", group: "Publish", color: "0", loc: "355 35.52284749830794" },
      { key: 9, text: "text for twoE", group: "Reviewing", color: "0", loc: "463 35.52284749830794" },
      { key: 10, text: "text for twoF", group: "Reviewing", color: "1", loc: "463 65.52284749830794" },
      { key: 11, text: "text for twoG", group: "Testing", color: "0", loc: "574 35.52284749830794" },
      { key: 12, text: "text for fourA", group: "Customer", color: "1", loc: "683 35.52284749830794" },
      { key: 13, text: "text for fourB", group: "Customer", color: "1", loc: "683 65.52284749830794" },
      { key: 14, text: "text for fourC", group: "Customer", color: "1", loc: "683 95.52284749830794" },
      { key: 15, text: "text for fourD", group: "Customer", color: "0", loc: "683 125.52284749830794" },
      { key: 16, text: "text for fiveA", group: "Customer", color: "0", loc: "683 155.52284749830795" }
    ]
    myDiagram.model = myModel
  }
  
  onMounted(() => {
    init();
  });
  </script>
    
    