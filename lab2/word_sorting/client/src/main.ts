import "./style.css";

const socket = new WebSocket("ws://localhost:8080/sort");

const input = document.querySelector<HTMLTextAreaElement>("#input")!;

const resultContainer = document.querySelector<HTMLDivElement>("#result")!;

document
  .querySelector<HTMLButtonElement>("#sort")!
  .addEventListener("click", async () => {
    socket.send(input!.value)
  });


socket.onmessage = (event) => {
  const output = event.data!;
  let words = output.split(" ");

  console.log("Recieved from server:")
  console.log(output)
  
  resultContainer.innerHTML = `
    <ul class='list-item'>
      ${words.map((word: string) => `<li>${word}</li>`).join("")}
    </ul>
  `
}
