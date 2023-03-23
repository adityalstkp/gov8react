import { linkStyle, wordStyle } from "../../styles/common";

const About = () => {
  return (
    <div>
      <h1 className={wordStyle}>React SSR with Go V8 Binding</h1>
      <a className={linkStyle} href="/">
        Home
      </a>
    </div>
  );
};

export default About;
