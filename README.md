## Getting Started

To use the TodoMVC framework, follow these steps:

1. **Import the Framework**

   Import the `FrameWork` and `Router` classes from the framework files into your project:

   ```javascript
   import { FrameWork } from "./framework/main";
   import Router from "./framework/router";
   ```

2. **Define Your Main Component**

   Create a main function that represents one route. This function will define your application's structure and behavior.

   ```
   function todoMVC(){
       main page ...
   }
   ```

3. **Use State Management**

   Use the FrameWork.UseState function to manage the state of your application. This function takes three arguments: the initial state value, a unique key to identify the state, and the parent component

   ```
   const [todos, setTodos] = FrameWork.UseState([], "todos", TodoMVC);
   ```

4. **Handle User Interactions**

   Implement functions to handle user interactions, such as adding, toggling, and removing elements. These functions should update the state using the `setState` function provided by `FrameWork` in `UseState`
   `setState` is `setTodos` in the example

   ```
   const addTodo = (event) => {
        // Add todo logic
        setTodos([...todos, { text: newValue, completed: false }]);
    };
   ```

5. **Create Elements**

   Use the `FrameWork.CreateElement` function to create and render HTML elements. This function takes three arguments: the element type, attributes, and child elements like text or other CreateElements.

   In the example there is a normal CreateElement, 2 nested elements, first has multible attributes, second has an event as an argument

   ```
    FrameWork.CreateElement("div", { className: "my-class" },
        FrameWork.CreateElement("input", {
            type: "text",
            value: variable,
            placeholder: "Add a new todo",
        }),
        FrameWork.CreateElement("button", { onClick: function() }, "button")
    );
   ```

6. **Router Setup**

   Initialize a `Router` instance to handle routing in your application. This is useful for creating different views or pages within your app.

   ```
   const container = document.getElementById("app");
   const router = new Router(container);
   ```

7. **Register Routes**

   Use the `registerRoute()` method to define routes and associate them with components or views. In the provided code, a default route / is registered to render the page.

   ```
    router.registerRoute('/', () => TodoMVC());
   ```

8. **Handle Route Changes**

   Finally, call `handleRouteChange()` to handle initial route setup and respond to route changes.

   ```
    router.handleRouteChange();
   ```
