export function DiseaseIcon({ className }: { className?: string }) {
  return (
    <svg
      className={className}
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth="2"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
      />
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M16 12a4 4 0 11-8 0 4 4 0 018 0z"
      />
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        d="M12 12v.01"
      />
    </svg>
  );
}