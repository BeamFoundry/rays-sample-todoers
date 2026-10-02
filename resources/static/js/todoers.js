var lists = null;
var items = []

function addItemToForm(form, name, value) {
    let input = document.createElement('input');
    input.type = "hidden";
    input.name = name;
    input.value = value;
    form.append(input);
}

document.addEventListener("htmx:responseError", (event) => {
    const xhr = event.detail.xhr
    if (xhr.status != 200) {
        event.stopPropagation(); // Tell htmx not to process these requests
        location.reload();
    }
})
document.addEventListener("DOMContentLoaded", (event) => {
    setTimeout(function() {
        document.querySelectorAll(".alert").forEach(function (alert) {
            new bootstrap.Alert(alert).close();
        });
    }, 3000);
});
htmx.onLoad(function(content) {
    var container = document.querySelector("#todolists-container");
    if (lists !== null) {
        lists.destroy();
    }
    while (items.length > 0) {
        var item = items.pop();
        item.destroy();
    }
    lists = new Sortable(container, {
        onEnd: function (evt) {
            let form = document.querySelector('#rank-todolist');
            form.replaceChildren();
            let todos = container.querySelectorAll(".todolist");
            todos.forEach(function (todo, index, list) {
                let id = todo.id.replace(/todolist-/,"");
                addItemToForm(form, "todo", id);
            });
            htmx.trigger(form, "submit");
        }
    });
    var todolists = container.querySelectorAll(".todolist");
    for (var i = 0; i < todolists.length; i++) {
        let todolist = todolists[i];
        let item_container = todolist.querySelector(".items-container");
        let form = todolist.querySelector("#actions-"+todolist.id);
        //console.log(item_container)
        var item = new Sortable(item_container, {
            animation: 150,
            group: "items",
            onEnd: function (evt) {
                if (evt.from == evt.to) {
                    form.replaceChildren();
                    addItemToForm(form, "action", "reorder");
                    let list_items = todolist.querySelectorAll(".todolist-item");
                    list_items.forEach(function (list_item, index, list) {
                        let id = list_item.id.replace(/item-[0-9]+-/,"");
                        addItemToForm(form, "data", id);
                    });
                    htmx.trigger(form, "submit");
                } else {
                    let id = evt.to.id.replace("items-container-", "");
                    form.replaceChildren();
                    addItemToForm(form, "action", "move");
                    addItemToForm(form, "data", evt.oldIndex + 1);
                    addItemToForm(form, "data", id);
                    addItemToForm(form, "data", evt.newIndex + 1);
                    htmx.trigger(form, "submit");
                }
            }
        });
        items.push(item);
    }
});
