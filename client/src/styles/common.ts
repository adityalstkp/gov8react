import { css } from "@emotion/css";

export const wordStyle = css`
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  padding: 4rem 0;
  font-weight: bold;
  font-family: var(--font-mono);
  text-align: center;
  font-style: normal;
  letter-spacing: 1px;
`;

export const subTitleStyle = css`
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  font-weight: bold;
  font-family: var(--font-mono);
  text-align: center;
  font-style: normal;
  letter-spacing: 1px;
`;

export const linkStyle = css`
  display: flex;
  justify-content: center;
  align-items: center;
  text-decoration: none;
  color: inherit;
  font-family: var(--font-mono);
  text-align: center;
`;

export const containerStyle = css`
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  margin: auto;
  height: 200px;
  width: 500px;

  & {
    :before,
    ::after {
      content: "";
      left: 50%;
      position: absolute;
      filter: blur(45px);
      transform: translateZ(0);
    }

    ::before {
      background: var(--primary-glow);
      border-radius: 50%;
      width: 480px;
      height: 360px;
      margin-left: -400px;
    }

    ::after {
      background: var(--secondary-glow);
      width: 240px;
      height: 180px;
      z-index: -1;
    }
  }
`;
