import type { ReactNode } from 'react';

export function PageFrame({
  eyebrow,
  title,
  subtitle,
  aside,
  children,
}: {
  eyebrow: string;
  title: string;
  subtitle?: string;
  aside?: ReactNode;
  children: ReactNode;
}) {
  return (
    <div className="page-frame">
      <div className="page-hero">
        <div>
          <div className="eyebrow">{eyebrow}</div>
          <h1>{title}</h1>
          {subtitle ? <p className="hero-copy">{subtitle}</p> : null}
        </div>
        {aside ? <div className="hero-aside">{aside}</div> : null}
      </div>
      {children}
    </div>
  );
}
