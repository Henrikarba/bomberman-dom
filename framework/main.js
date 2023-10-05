let hooks = {}
let globalContainer;

function InitFramework(element, container) {
  globalContainer = container
  Render(element, container)
}

function createTextElement(text) {
  return {
    type: "TEXT_ELEMENT",
    props: {
      nodeValue: text,
      children: [],
    },
  }
}

function CreateElement(type, props, ...children) {
  props = props || {} // add this line to ensure props is not null or undefined
  const childElements = Array.isArray(children)
    ? children.flatMap(child =>
      typeof child === "object"
        ? child
        : createTextElement(child)
    )
    : typeof children === "object"
      ? [children]
      : [createTextElement(children)]
  const element = {
    type,
    props: {
      ...props,
      children: childElements,
    },
  }
  return element
}
function isEvent(propName) {
  return propName.startsWith("on")
}

function Render(element, container) {
  const dom =
    element.type == "TEXT_ELEMENT"
      ? document.createTextNode("")
      : document.createElement(element.type)
  const isProperty = key => key !== "children"
  Object.keys(element.props)
    .filter(isProperty)
    .forEach(name => {
      if (isEvent(name)) {
        const eventType = name.toLowerCase().substring(2)
        dom.addEventListener(eventType, element.props[name])
      } else {
        dom[name] = element.props[name]
      }
    })
  element.props.children.forEach(child =>
    Render(child, dom)
  )
  container.appendChild(dom)
}

function UseState(initialValue, key) {
  console.log("hooks", hooks)
  if (!hooks.hasOwnProperty(key)){
    hooks[key] = initialValue
  }
  const setState = (newValue, element, key) => {
    hooks[key] = newValue;
      globalContainer.innerHTML = "";
      Render(element(), globalContainer)
  };
  return [hooks[key], setState];
}

const FrameWork = {
  CreateElement,
  Render,
  UseState,
  InitFramework,
}

export { FrameWork }