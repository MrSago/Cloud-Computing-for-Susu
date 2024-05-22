const addForm = document.getElementById("add-form");

const editFormSection = document.getElementById("edit-form-section");
const editForm = document.getElementById("edit-form");
const editFormNameField = document.getElementById("edit-name");
const editFormIdField = document.getElementById("edit-id");
const editFormIsCompleteField = document.getElementById("edit-isComplete");
const editCloseButton = document.getElementById("edit-close");

let todos = [];

addForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const addNameTextbox = document.getElementById("add-name");
  const name = addNameTextbox.value.trim();
  addNameTextbox.value = "";

  await fetch("api/todoitems", {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      isComplete: false,
      name: name,
    }),
  });

  updateItemsView();
});

editForm.addEventListener("submit", async (event) => {
  event.preventDefault();

  const itemId = parseInt(editFormIdField.value);
  const isItemComplete = editFormIsCompleteField.checked;
  const itemName = editFormNameField.value.trim();

  await fetch(`api/todoitems/${itemId}`, {
    method: "PUT",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      id: itemId,
      isComplete: isItemComplete,
      name: itemName,
    }),
  });

  editFormSection.style.display = "none";

  await updateItemsView();
});

editCloseButton.addEventListener("click", (event) => {
  event.preventDefault();

  editFormSection.style.display = "none";
});

async function deleteItem(id) {
  await fetch(`api/todoitems/${id}`, {
    method: "DELETE",
  });

  await updateItemsView();
}

function displayEditForm(id) {
  const item = todos.find((item) => item.id === id);
  if (!item) {
    return;
  }

  console.log("displayEditForm", editFormSection);

  editFormSection.style.display = "block";

  editFormNameField.value = item.name;
  editFormIdField.value = item.id;
  editFormIsCompleteField.checked = item.isComplete;
}

function _displayCount(itemCount) {
  const name = itemCount === 1 ? "to-do" : "to-dos";
  document.getElementById("counter").innerText = `${itemCount} ${name}`;
}

function _displayItems(items) {
  _displayCount(items.length);

  const tBody = document.getElementById("todos");
  tBody.innerHTML = "";

  const button = document.createElement("button");

  items.forEach((item) => {
    let isCompleteCheckbox = document.createElement("input");
    isCompleteCheckbox.type = "checkbox";
    isCompleteCheckbox.disabled = true;
    isCompleteCheckbox.checked = item.isComplete;

    let editButton = button.cloneNode(false);
    editButton.innerText = "Edit";
    editButton.addEventListener("click", () => displayEditForm(item.id));

    let deleteButton = button.cloneNode(false);
    deleteButton.innerText = "Delete";
    deleteButton.addEventListener("click", () => deleteItem(item.id));

    let tr = tBody.insertRow();

    let td1 = tr.insertCell(0);
    td1.appendChild(isCompleteCheckbox);

    let td2 = tr.insertCell(1);
    let textNode = document.createTextNode(item.name);
    td2.appendChild(textNode);

    let td3 = tr.insertCell(2);
    td3.appendChild(editButton);

    let td4 = tr.insertCell(3);
    td4.appendChild(deleteButton);
  });

  todos = items;
}

async function updateItemsView() {
  const res = await fetch("api/todoitems");
  const items = await res.json();
  _displayItems(items);
}

updateItemsView();
