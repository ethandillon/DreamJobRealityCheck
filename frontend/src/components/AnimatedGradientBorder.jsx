// src/components/AnimatedGradientBorder.jsx
import React, { useRef, useEffect, CSSProperties } from 'react';

// We'll rename it slightly for clarity in our project
export const AnimatedGradientBorder = ({ children, className }) => {
  const boxRef = useRef(null);

  useEffect(() => {
    const boxElement = boxRef.current;
    if (!boxElement) return;

    const updateAnimation = () => {
      const angle = (parseFloat(boxElement.style.getPropertyValue("--angle")) + 0.5) % 360;
      boxElement.style.setProperty("--angle", `${angle}deg`);
      requestAnimationFrame(updateAnimation);
    };
    requestAnimationFrame(updateAnimation);
  }, []);

  // We are defining the base styles here, but allowing overrides via className prop
  const defaultClassName = "border-2 border-[#0000] [background:padding-box_var(--bg-color),border-box_var(--border-color)]";

  return (
    <div
      ref={boxRef}
      style={
        {
          "--angle": "0deg",
          // We'll make the gradient a bit more vibrant to match our theme
          "--border-color": "linear-gradient(var(--angle), #ffffffff, #283a5cff, #ffffffff)", 
          // Default background color from the original component
          "--bg-color": "linear-gradient(#131219, #131219)", 
        }
      }
      className={`${defaultClassName} ${className}`} // Combine default and passed classNames
    >
      {children}
    </div>
  );
};