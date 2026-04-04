/*
Copyright (C) 2025 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/

import React from 'react';

const Contact = () => {
  return (
    <div className='mt-[60px] flex min-h-[calc(100vh-60px)] flex-col items-center justify-center px-4 py-12'>
      <h1 className='mb-2 text-3xl font-bold text-gray-900 dark:text-white'>
        请添加我们的微信
      </h1>
      <p className='mb-8 text-base text-gray-500 dark:text-gray-400'>
        扫描下方二维码，添加微信
      </p>

      <div
        className='flex items-center justify-center rounded-2xl p-8'
        style={{
          background: 'linear-gradient(135deg, #1a7fe8 0%, #0d6ed4 60%, #4db8f8 100%)',
          boxShadow: '0 8px 32px rgba(13, 110, 212, 0.35)',
          minWidth: 320,
          minHeight: 380,
        }}
      >
        {/* WeChat bubble decoration */}
        <div className='relative flex flex-col items-center'>
          {/* Circular glow behind QR */}
          <div
            className='absolute'
            style={{
              width: 280,
              height: 280,
              borderRadius: '50%',
              background: 'rgba(77, 184, 248, 0.25)',
              top: '50%',
              left: '50%',
              transform: 'translate(-50%, -50%)',
            }}
          />

          {/* QR Code card */}
          <div
            className='relative z-10 rounded-2xl bg-white p-4'
            style={{ boxShadow: '0 4px 20px rgba(0,0,0,0.15)' }}
          >
            <img
              src='/wechat-qrcode.png'
              alt='微信二维码'
              style={{ width: 220, height: 220, display: 'block' }}
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Contact;
