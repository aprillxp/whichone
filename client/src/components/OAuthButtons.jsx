export default function oAuthButton() {
  const API = import.meta.env.BASE_URL;
  return (
    <div className="flex flex-col gap-2">
      <a className="btn btn-outline" href={`${API}/auth/google`}>
        Google Login
      </a>
      <a className="btn btn-outline" href={`${API}/auth/twitter`}>
        Twitter Login
      </a>
    </div>
  );
}
