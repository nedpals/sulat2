import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { isEqual, cloneDeep } from "lodash"
import { useRef } from "react";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function useDeepCompareMemoize<T>(value: T, defaultValue: T) {
  const ref = useRef<T>();
  if (!isEqual(value, ref.current)) {
    ref.current = cloneDeep(value);
  }
  return ref.current ?? defaultValue;
};
