import "./style.css";

document.querySelector<HTMLDivElement>("#app")!.innerHTML = `
  <div class="chat-container">
    <div class="messages" id="messages">
        <ul class="item-ul">
            <!-- Messages will be added dynamically here -->
        </ul>
    </div>
    <div class="input-container">
        <input type="text" id="userInput" placeholder="Type your message...">
        <input type="submit" id="submitButton" value="Send" onclick="sendMessage()">
    </div>
  </div>
`;

const socket = new WebSocket("ws://localhost:8080/bulletin_board");

const input = document.querySelector<HTMLInputElement>("#userInput")!;
const submitButton = document.querySelector<HTMLInputElement>("#submitButton")!;

input.addEventListener("keypress", function (event) {
  if (event.key === "Enter") {
    event.preventDefault();
    submitButton.click();
  }
});

submitButton.onclick = () => {
  const message = input.value;
  socket.send(JSON.stringify(message));
  input.value = "";
};

type ServerResponse = {
  AnswerType: string;
  Value: string[];
};

socket.onmessage = (event) => {
  const output: ServerResponse = JSON.parse(event.data);
  if (output.AnswerType === "LIST") {
    let messages = output.Value;
    let messageList = document.querySelector<HTMLUListElement>(".item-ul")!;
    if (messages.length == 0) {
      messageList.innerHTML = `
        <h2>No items</h2>
        `;
    } else {
      messageList.innerHTML = `
      <h2>Items</h2>
      ${messages.map((message: string) => `<li>${message};</li>`).join("")}
    `;
    }
  } else if (output.AnswerType === "MESSAGE") {
    let messageList = document.querySelector<HTMLUListElement>(".item-ul")!;
    messageList.innerHTML = `<h3>Message</h3>\n<li>${output.Value}</li>`;
  }
};
