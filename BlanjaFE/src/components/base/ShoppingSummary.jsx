import React from "react";

const CartSummary = ({ total, handleClick }) => {
  const formatRupiah = (price) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(price);
  };
  return (
    <div className="w-full max-lg:h-40 rounded-lg max-lg:rounded-lg p-5 shadow-[0_0_14px_0_#ADADAD40]">
      <div className=" mt-2 max-lg:hidden">
        <span className="text-lg font-bold text-[#222222]">
          Shopping Summary
        </span>
      </div>
      <div className="mt-8 max-lg:flex-col max-lg:flex max-lg:justify-between max-lg:h-full max-lg:mt-0">
        <div className="flex justify-between max-lg:h-1/2">
          <span className="text-lg font-bold text-gray-300">Total Price</span>
          <span className="text-lg font-bold">{formatRupiah(total)}</span>
        </div>

        <div className="max-lg:h-1/2">
          <button onClick={handleClick} className="w-full bg-red-500 text-white px-4 py-2 rounded-full mt-8 max-lg:mt-0">
            Check Out
          </button>
        </div>
      </div>
    </div>
  );
};

export default CartSummary;
