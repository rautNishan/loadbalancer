import express from "express";
const app = express();
const port = 3002;

app.listen(port, () => {
  console.log(`Server running on http://localhost:${port}`);
});

app.get("/", (req, res) => {
  console.log("Incoming Request In server 2");
  res.send("Hello world");
});

app.get("/test", (req, res) => {
  console.log("Incoming Request In server 2");
  for (let i = 0; i < 1000000; i++) {}
  res.send("Loop World");
});
