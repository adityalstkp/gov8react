import { useNavigate } from "react-router-dom";
import { useIntro } from "../../hooks/intro";
import { containerStyle, linkStyle, wordStyle } from "../../styles/common";

const Home = () => {
  const { data: greet } = useIntro();
  const navigate = useNavigate();

  const handleToAbout = () => {
    navigate("/about");
  };

  return (
    <div className={containerStyle}>
      <h1 className={wordStyle}>gov8react</h1>
      {greet && <p className={wordStyle}>{greet}</p>}
      <a className={linkStyle} onClick={handleToAbout}>
        About
      </a>
    </div>
  );
};

export default Home;
