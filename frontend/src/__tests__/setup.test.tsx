import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: { retry: false },
  },
});

describe('Test Setup', () => {
  it('should render test components', () => {
    const TestComponent = () => <div>Test Component</div>;
    
    render(
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <TestComponent />
        </BrowserRouter>
      </QueryClientProvider>
    );
    
    expect(screen.getByText('Test Component')).toBeDefined();
  });

  it('should handle React Query', () => {
    expect(queryClient).toBeDefined();
    expect(queryClient.getQueryCache()).toBeDefined();
  });
});

