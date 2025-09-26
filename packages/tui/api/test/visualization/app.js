// Fallback data in case fetch fails
const fallbackData = {
  nodes: [
    { name: "internal.Session" },
    { name: "internal.Message" },
    { name: "internal.Part" },
    { name: "internal.Config" },
    { name: "internal.Provider" },
    { name: "internal.FileContent" },
    { name: "internal.Symbol" },
    { name: "internal.Auth" },
    { name: "internal.Event" },
    { name: "internal.Command" },
    { name: "internal.Permission" },
    { name: "internal.Error" },
    { name: "Wrap in Opt*" },
    { name: "Add Discriminator" },
    { name: "Apply Validation" },
    { name: "Nest Structures" },
    { name: "Build Sum Type" },
    { name: "Session" },
    { name: "Message" },
    { name: "Part" },
    { name: "Config" },
    { name: "Provider" },
    { name: "FileContent" },
    { name: "Symbol" },
    { name: "Auth" },
    { name: "Event" },
    { name: "Command" },
    { name: "Permission" },
    { name: "Error" },
  ],
  links: [
    { source: "internal.Session", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Session", value: 1, label: "Build Sum Type" },
    { source: "internal.Message", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Message", value: 1, label: "Build Sum Type" },
    { source: "internal.Part", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Part", value: 1, label: "Build Sum Type" },
    { source: "internal.Config", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Config", value: 1, label: "Build Sum Type" },
    { source: "internal.Provider", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Provider", value: 1, label: "Build Sum Type" },
    { source: "internal.FileContent", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "FileContent", value: 1, label: "Build Sum Type" },
    { source: "internal.Symbol", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Symbol", value: 1, label: "Build Sum Type" },
    { source: "internal.Auth", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Auth", value: 1, label: "Build Sum Type" },
    { source: "internal.Event", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Event", value: 1, label: "Build Sum Type" },
    { source: "internal.Command", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Command", value: 1, label: "Build Sum Type" },
    { source: "internal.Permission", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Permission", value: 1, label: "Build Sum Type" },
    { source: "internal.Error", target: "Wrap in Opt*", value: 1, label: "Wrap in Opt*" },
    { source: "Wrap in Opt*", target: "Add Discriminator", value: 1, label: "Add Discriminator" },
    { source: "Add Discriminator", target: "Apply Validation", value: 1, label: "Apply Validation" },
    { source: "Apply Validation", target: "Nest Structures", value: 1, label: "Nest Structures" },
    { source: "Nest Structures", target: "Build Sum Type", value: 1, label: "Build Sum Type" },
    { source: "Build Sum Type", target: "Error", value: 1, label: "Build Sum Type" },
  ],
}

let data
try {
  const response = await fetch("data.json")
  data = await response.json()
} catch (e) {
  console.info("Using fallback data due to fetch/module constraints.")
  data = fallbackData
}

const tooltipTemplate = Handlebars.compile("<strong>{{name}}</strong>")

function renderTooltip(node) {
  return tooltipTemplate(node)
}

function buildSankey(data) {
  const container = d3.select("#sankey-container")
  const svg = container.select("svg")
  const width = svg.node().getBoundingClientRect().width
  const height = svg.node().getBoundingClientRect().height

  const sankey = d3
    .sankey()
    .nodeWidth(15)
    .nodePadding(20)
    .nodeId((d) => d.name)
    .extent([
      [10, 10],
      [width - 10, height - 10],
    ])

  const graph = sankey({
    nodes: data.nodes.map((d) => ({ ...d })),
    links: data.links.map((d) => ({ ...d })),
  })

  // Create a group for zoom
  const g = svg.append("g").attr("class", "sankey-group")

  // Links
  g.append("g")
    .attr("class", "links")
    .selectAll("path")
    .data(graph.links)
    .join("path")
    .attr("d", d3.sankeyLinkHorizontal())
    .attr("stroke-width", (d) => Math.max(1, d.width))
    .classed("link", true)

  // Nodes
  const node = g
    .append("g")
    .attr("class", "nodes")
    .selectAll("g")
    .data(graph.nodes)
    .join("g")
    .attr("class", "node")
    .attr("tabindex", 0)
    .attr("aria-label", (d) => `Node: ${d.name}`)

  node
    .append("rect")
    .attr("x", (d) => d.x0)
    .attr("y", (d) => d.y0)
    .attr("height", (d) => d.y1 - d.y0)
    .attr("width", (d) => d.x1 - d.x0)

  node
    .append("text")
    .attr("x", (d) => d.x0 - 6)
    .attr("y", (d) => (d.y1 + d.y0) / 2)
    .attr("dy", "0.35em")
    .attr("text-anchor", "end")
    .text((d) => d.name)
    .filter((d) => d.x0 < width / 2)
    .attr("x", (d) => d.x1 + 6)
    .attr("text-anchor", "start")

  // Zoom
  const zoom = d3
    .zoom()
    .scaleExtent([0.1, 10])
    .on("zoom", (event) => {
      g.attr("transform", event.transform)
    })

  svg.call(zoom)

  // Drag
  node.call(
    d3
      .drag()
      .on("start", function () {
        d3.select(this).raise()
      })
      .on("drag", function (event, d) {
        d.x0 += event.dx
        d.x1 += event.dx
        d.y0 += event.dy
        d.y1 += event.dy
        d3.select(this).select("rect").attr("x", d.x0).attr("y", d.y0)
        d3.select(this)
          .select("text")
          .attr("x", d.x0 - 6)
          .attr("y", (d.y1 + d.y0) / 2)
        // Update links
        g.selectAll(".link").attr("d", d3.sankeyLinkHorizontal())
      }),
  )

  return { g, graph }
}

function bindInteractions({ g, graph }) {
  const tooltip = d3.select("#tooltip")
  let highlighted = false

  function showTooltip(event, d) {
    tooltip
      .style("left", event.pageX + 10 + "px")
      .style("top", event.pageY - 10 + "px")
      .html(renderTooltip(d))
      .style("opacity", 1)
  }

  function hideTooltip() {
    tooltip.style("opacity", 0)
  }

  function highlightConnected(d) {
    const connectedNodes = new Set()
    const connectedLinks = new Set()

    // Traverse outgoing
    function traverse(node) {
      if (connectedNodes.has(node)) return
      connectedNodes.add(node)
      graph.links.forEach((link) => {
        if (link.source === node) {
          connectedLinks.add(link)
          traverse(link.target)
        }
        if (link.target === node) {
          connectedLinks.add(link)
          traverse(link.source)
        }
      })
    }
    traverse(d)

    g.classed("highlighted", true)
    g.selectAll(".node").classed("highlight", (n) => connectedNodes.has(n))
    g.selectAll(".link").classed("highlight", (l) => connectedLinks.has(l))
  }

  function resetHighlight() {
    g.classed("highlighted", false)
    g.selectAll(".node").classed("highlight", false)
    g.selectAll(".link").classed("highlight", false)
  }

  g.selectAll(".node")
    .on("mouseover", showTooltip)
    .on("mouseout", hideTooltip)
    .on("click", function (event, d) {
      event.stopPropagation()
      if (highlighted) {
        resetHighlight()
        highlighted = false
      } else {
        highlightConnected(d)
        highlighted = true
      }
    })
    .on("keydown", function (event, d) {
      if (event.key === "Enter") {
        event.preventDefault()
        if (highlighted) {
          resetHighlight()
          highlighted = false
        } else {
          highlightConnected(d)
          highlighted = true
        }
      }
    })

  g.selectAll(".link")
    .on("mouseover", function (event, d) {
      showTooltip(event, { name: d.label })
    })
    .on("mouseout", hideTooltip)

  // Click outside to reset
  d3.select("#sankey-container").on("click", function (event) {
    if (!event.target.closest(".node") && !event.target.closest(".link")) {
      resetHighlight()
      highlighted = false
    }
  })
}

const { g, graph } = buildSankey(data)
bindInteractions({ g, graph })
