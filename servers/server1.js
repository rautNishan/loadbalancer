import express from "express";
const app = express();
const port = 3001;

app.listen(port, () => {
  console.log(`Server running on http://localhost:${port}`);
});

app.get("/", (req, res) => {
  console.log("Incoming Request In server 1");
  res.send("Hello world");
});

app.get("/test", (req, res) => {
  console.log("Incoming Request In server 2");
  for (let i = 0; i < 9000000000; i++) {}
  res.send("Loop World");
});
