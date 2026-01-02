import { Link } from 'react-router-dom';

export default function LandingPage() {
  return (
    <div className="flex flex-col items-center justify-center text-center gap-6">
      <h1 className="text-4xl font-extrabold text-base-content">Choose smart, play better.</h1>

      <Link to="/register" className="btn btn-primary">
        Get Started
      </Link>
    </div>
  );
}
