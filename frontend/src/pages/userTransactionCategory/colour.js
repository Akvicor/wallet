export const RandomColour16 = () => {
  let H = Math.random();
  let S = Math.random();
  let L = Math.random();

  S = 0.7 + (S * 0.2); // [0.7 - 0.9] 排除过灰颜色
  L = 0.4 + (L * 0.4); // [0.4 - 0.8] 排除过亮过暗色

  H = parseFloat(H.toFixed(2))
  S = parseFloat(S.toFixed(2))
  L = parseFloat(L.toFixed(2))

  let R;
  let G;
  let B;
  if (+S === 0) {
    R = G = B = L; // 饱和度为0 为灰色
  } else {
    let hue2Rgb = function (p, q, t) {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1/6) return p + (q - p) * 6 * t;
      if (t < 1/2) return q;
      if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
      return p;
    };
    let Q = L < 0.5 ? L * (1 + S) : L + S - L * S;
    let P = 2 * L - Q;
    R = hue2Rgb(P, Q, H + 1/3);
    G = hue2Rgb(P, Q, H);
    B = hue2Rgb(P, Q, H - 1/3);
  }

  R = Math.round(R * 255).toString(16)
  G = Math.round(G * 255).toString(16)
  B = Math.round(B * 255).toString(16)

  if (R.length === 1) {
    R = '0'+R
  }
  if (G.length === 1) {
    G = '0'+G
  }
  if (B.length === 1) {
    B = '0'+B
  }

  return '#'+R+G+B;
}
