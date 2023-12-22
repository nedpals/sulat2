import Scaffold from "../components/Scaffold";

export default function RecordEditor() {
  return (
    <Scaffold
      className="flex"
      leftHeader={() => <div>
          <p className="text-sm text-slate-600 block">Collection Name</p>
          <h1 className="text-lg font-bold text-slate-800">Edit Record</h1>
      </div>}
      actions={() => <>
        <p className="text-sm italic text-slate-600">Last modified 8 minutes ago</p>
        <button className="sulat-btn bg-red-100/30 hover:bg-red-100 active:bg-red-200 text-red-500">Delete</button>
        <button className="sulat-btn is-primary px-12">Save</button>
      </>}>
      <div className="w-3/4 space-y-3 pr-2">
        <input
          type="text"
          className="sulat-input text-xl w-full"
          placeholder="Title" />

        <textarea className="bg-white rounded border w-full h-48">

        </textarea>
      </div>

      <div className="flex flex-col items-stretch w-1/4 pl-2 space-y-3">
        <div className="flex flex-col items-stretch space-y-2">
          <button className="sulat-btn">Preview</button>
        </div>

        <div className="border rounded bg-white shadow">
          <div className="px-3 py-2 flex items-center justify-between">
            <h3 className="font-medium">Categories</h3>
          </div>

          <div className="px-3 pb-2">
            <p>Test</p>
          </div>
        </div>
      </div>
    </Scaffold>
  );
}
