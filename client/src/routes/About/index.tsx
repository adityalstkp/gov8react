import { useNavigate } from "react-router-dom";
import { linkStyle, wordStyle } from "../../styles/common";

const About = () => {
  const navigate = useNavigate();
  const handleToHome = () => {
    navigate("/");
  };

  return (
    <div>
      <h1 className={wordStyle}>React SSR with Go V8 Binding</h1>
      <a className={linkStyle} onClick={handleToHome}>
        Home
      </a>
    </div>
  );
};

export default About;
