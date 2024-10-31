export interface QueryProblemListParams {
    page: number;
    page_size: number;
    keyword: string;
    tag: string;
}

export interface SubmitForm {
    problem_id: number;
    title: string;
    lang: string;
    code: string;
}

export interface Problem {
    id: number;
    title: string;
    description: string;
    level: number;
    tags: string[];
    status?: number;
    create_at?: string;
    create_by?: number;
    sample_input?: string;
    sample_output?: string;
    hint?: string;
}

export interface ProblemConfig {
    testCases: TestCase[];
    compileLimit: Limit;
    runLimit: Limit;
}

interface TestCase {
    input: string;
    output: string;
}

interface Limit {
    cpuLimit: number;
    clockLimit: number;
    memoryLimit: number;
    stackLimit: number;
    procLimit: number;
}